'''
Yandex Cloud
'''

from dataclasses import dataclass

import pulumi_yandex as yandex

from . import timestamp
from .cloud import CloudProvider, CloudInstance, status
from .logging import log
from .templating import template


ROUTER_IP   = '10.10.10.10'
ROUTER_CIDR = '10.10.10.0/24'
INNER_CIDR  = '10.0.0.0/24'


@dataclass(eq=False)
class YandexInstance(CloudInstance):
    '''Yandex Cloud VPS'''

    ipv4_address: str = None

    def update_status(self):
        '''Write updated values to self.status, self.idle_since'''
        super().update_status()

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
                    gitlab_runner_token = 'gitlab_runner_token', # TODO
                ),
                'serial-port-enable': 1,
            },
            zone = config.availability_zone,
            resources = yandex.ComputeInstanceResourcesArgs(
                cores = config.vcpu_count,
                memory = config.memory_gb,
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
        super().create()

    def cleanup(self):
        '''
        Prepare instance for deletion:
            - Unregister GitLab runners
            - Remove cloud firewall rules
            - etc.

        After successful cleanup this method must set self.status to DESTROYING
        '''
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
        self.vpc = yandex.VpcNetwork(f'{__package__}:{self.__class__.__name__}:network')
        router_table = yandex.VpcRouteTable(
                'nat',
                network_id=self.vpc.id,
                static_routes=[yandex.VpcRouteTableStaticRouteArgs(
                    destination_prefix='0.0.0.0/0',
                    next_hop_address=ROUTER_IP,
                )],
        )
        router_subnet = yandex.VpcSubnet(
                f'{__package__}:{self.__class__.__name__}:router',
                network_id = self.vpc.id,
                zone=config.availability_zone,
                v4_cidr_blocks=[ROUTER_CIDR],
        )
        self.subnet = yandex.VpcSubnet(
                f'{__package__}:{self.__class__.__name__}:subnet',
                network_id = self.vpc.id,
                zone=config.availability_zone,
                v4_cidr_blocks=[INNER_CIDR],
                route_table_id=router_table.id,
        )
        self.ipaddr = yandex.VpcAddress(
                'addr',
                external_ipv4_address=yandex.VpcAddressExternalIpv4AddressArgs(
                    zone_id=config.availability_zone,
                )
        )
        self.router = yandex.ComputeInstance(
            resource_name = 'router',
            hostname = 'router',
            scheduling_policy = yandex.ComputeInstanceSchedulingPolicyArgs(
                preemptible = config.preemptible_instances,
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
            resources = yandex.ComputeInstanceResourcesArgs(
                cores = 2,
                memory = 2,
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
                    nat_ip_address=self.ipaddr.external_ipv4_address.address,
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
