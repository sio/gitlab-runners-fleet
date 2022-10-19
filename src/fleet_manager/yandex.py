'''
Yandex Cloud (https://yandex.cloud)

This module contains provider specific Pulumi code to launch runner instances.
Each runner (CloudInstance) provisions one compute instance from a clean Debian
image.

Extra billable resources (neccessary overhead):
    - One public static IP address
    - One compute instance for NAT router

Pricing depends on:
    - Number of resources assigned to VM (vCPU, RAM, storage)
    - CPU generation (aka platform in Yandex docs):
      https://cloud.yandex.com/en-ru/docs/compute/concepts/vm-platforms
    - Guaranteed CPU utilization (5%/20%/100%)

More information on prices:
    https://cloud.yandex.com/en-ru/docs/compute/pricing#prices
'''

from dataclasses import dataclass

import requests
import pulumi_yandex as yandex

from . import timestamp
from .cloud import CloudProvider, CloudInstance, status
from .logging import log
from .templating import template


ROUTER_IP   = '10.10.10.10'
ROUTER_CIDR = '10.10.10.0/24'
INNER_CIDR  = '10.0.0.0/24'


@dataclass(eq=False, repr=False)
class YandexInstance(CloudInstance):
    '''Yandex Cloud VPS'''

    ipv4_address: str = ''
    jobs_total: int = 0

    def update_status(self):
        '''Write updated values to self.status, self.idle_since'''
        self.status = self._fetch_status()
        super().update_status()

    def _fetch_status(self):  # TODO: move generic parts into CloudInstance
        scaling = self.cloud.scaling
        try:
            if not getattr(self, 'ipv4_address'):
                raise ValueError(f'ipv4_address not found for {self}')
            metrics_url = f'http://{self.ipv4_address}/metrics'
            log.debug('Fetching instance metrics: %s (%s)', self.name, metrics_url)
            response = requests.get(metrics_url, headers={'Host': self.name})
            response.raise_for_status()
            metrics = response.json()
        except Exception as exc:
            log.debug('Error while fetching metrics: %s', exc)
            if timestamp.now() - self.created_at < scaling.est_provisioning_minutes * 60:
                return status.PROVISIONING
            return status.ERROR
        jobs_total = metrics.get('gitlab_runner_jobs_total', 0)
        if jobs_total > self.jobs_total:
            self.jobs_total = jobs_total
            self.idle_since = 0
        if metrics.get('gitlab_runner_jobs', 0) > 0:
            return status.BUSY
        if self.idle_since and timestamp.now() - self.idle_since > scaling.max_idle_minutes * 60:
            return status.IDLE
        return status.READY

    def create(self):
        '''Create cloud instance corresponding to this object'''
        config = self.cloud.config
        instance = yandex.ComputeInstance(
            resource_name = self.name,
            hostname = self.name,
            scheduling_policy = yandex.ComputeInstanceSchedulingPolicyArgs(
                preemptible = config.preemptible_instances,
            ),
            metadata = {
                # https://cloud.yandex.com/en-ru/docs/compute/concepts/vm-metadata
                'user-data': template(config.cloudinit_runner).render(
                    public_key = config.public_key,
                    gitlab_runner_token = self.cloud.gitlab.runner_token,
                ),
                'serial-port-enable': 1,
            },
            zone = config.availability_zone,
            platform_id = config.cpu_platform,
            resources = yandex.ComputeInstanceResourcesArgs(
                cores = config.vcpu_count,
                memory = config.memory_gb,
                core_fraction = config.vcpu_performance_percent,
            ),
            boot_disk=yandex.ComputeInstanceBootDiskArgs(
                auto_delete=True,
                initialize_params=yandex.ComputeInstanceBootDiskInitializeParamsArgs(
                    size=config.disk_size_gb,
                    image_id=yandex.get_compute_image(family=config.image_family).image_id,
                ),
            ),
            network_interfaces=[yandex.ComputeInstanceNetworkInterfaceArgs(
                subnet_id=self.cloud.subnet.id,
            )],
        )
        if not getattr(self, 'ipv4_address'):
            self.ipv4_address = self.cloud.ipv4_external.apply(str)
        super().create()

    def cleanup(self):
        '''
        Prepare instance for deletion:
            - Unregister GitLab runners
            - Remove cloud firewall rules
            - etc.

        After successful cleanup this method must set self.status to DESTROYING
        '''
        log.debug(f'Initiating cleanup: {self}')
        try:
            if not getattr(self, 'ipv4_address'):
                raise ValueError(f'ipv4_address not found for {self}')
            response = requests.post(
                f'http://{self.ipv4_address}/unregister',
                headers={'Host': self.name},
            )
            response.raise_for_status()
        except Exception as exc:
            log.error(f'Cleanup failed: {exc}')
            self.status = status.ERROR
        super().cleanup()


class YandexCloud(CloudProvider):
    '''Yandex Cloud'''

    _instance_cls = YandexInstance

    def setup(self):
        '''
        Ensure that cloud provider is ready for creating instances:
            - Create required SSH keys
            - Configure cloud networking
            - Configure cloud NAT/firewall
            - etc.
        '''
        config = self.config
        vpc = yandex.VpcNetwork(f'{__package__}:{self.__class__.__name__}:network')
        router_table = yandex.VpcRouteTable(
                f'{__package__}:{self.__class__.__name__}:nat',
                network_id=vpc.id,
                static_routes=[yandex.VpcRouteTableStaticRouteArgs(
                    destination_prefix='0.0.0.0/0',
                    next_hop_address=ROUTER_IP,
                )],
        )
        router_subnet = yandex.VpcSubnet(
                f'{__package__}:{self.__class__.__name__}:router',
                network_id = vpc.id,
                zone=config.availability_zone,
                v4_cidr_blocks=[ROUTER_CIDR],
        )
        self.subnet = yandex.VpcSubnet(
                f'{__package__}:{self.__class__.__name__}:subnet',
                network_id = vpc.id,
                zone=config.availability_zone,
                v4_cidr_blocks=[INNER_CIDR],
                route_table_id=router_table.id,
        )
        ipv4_address = yandex.VpcAddress(
                f'{__package__}:{self.__class__.__name__}:address',
                external_ipv4_address=yandex.VpcAddressExternalIpv4AddressArgs(
                    zone_id=config.availability_zone,
                )
        )
        self.ipv4_external = ipv4_address.external_ipv4_address.address
        router = yandex.ComputeInstance(
            resource_name = 'router',
            hostname = 'router',
            scheduling_policy = yandex.ComputeInstanceSchedulingPolicyArgs(
                preemptible = False,
            ),
            metadata = {
                # https://cloud.yandex.com/en-ru/docs/compute/concepts/vm-metadata
                'user-data': template(config.cloudinit_router).render(
                    public_key = config.public_key,
                    inner_subnet = INNER_CIDR,
                ),
                'serial-port-enable': 1,
            },
            zone = config.availability_zone,
            platform_id = config.cpu_platform,
            resources = yandex.ComputeInstanceResourcesArgs(
                # router config is intentionally not tied to instance config
                cores = 2,
                memory = 2,
                core_fraction = 20,
            ),
            boot_disk=yandex.ComputeInstanceBootDiskArgs(
                auto_delete=True,
                initialize_params=yandex.ComputeInstanceBootDiskInitializeParamsArgs(
                    size=10,
                    image_id=yandex.get_compute_image(family=config.image_family).image_id,
                ),
            ),
            network_interfaces=[
                yandex.ComputeInstanceNetworkInterfaceArgs(
                    subnet_id=router_subnet.id,
                    nat=True,
                    nat_ip_address=self.ipv4_external,
                    ip_address=ROUTER_IP,
                ),
            ],
        )

    def _restore_from_deployment(self, stack):
        resources = stack.export_stack().deployment.get('resources', [])
        for resource in resources:
            if resource['type'] != 'yandex:index/computeInstance:ComputeInstance':
                continue
            if resource['outputs']['status'] != 'running':
                continue
            if resource['outputs']['hostname'] == 'router':
                continue
            params = dict(
                name=resource['outputs']['hostname'],
                created_at=timestamp.from_string(resource['outputs']['createdAt']),
            )
            instance = self._instance_cls(cloud=self, **params)
            self.instances.add(instance)
            log.debug(
                'Restored %i instances currently deployed: %s',
                len(self.instances),
                [i.name for i in self.instances]
            )
