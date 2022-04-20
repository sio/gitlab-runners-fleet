'''
Yandex Cloud
'''

from dataclasses import dataclass

import pulumi_yandex as yandex

from . import timestamp
from .cloud import CloudProvider, CloudInstance, status
from .logging import log
from .provision import template


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
                'user-data': template(config.cloudinit_template).render(
                    pubkey = 'publickey', # TODO
                    gitlab_runner_token = 'gitlab_runner_token', # TODO
                ),
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
        self.vpc = yandex.VpcNetwork(f'{__package__}:{self.__class__.__name__}:network')
        self.subnet = yandex.VpcSubnet(
                f'{__package__}:{self.__class__.__name__}:subnet',
                network_id = self.vpc.id,
                zone=self.config.availability_zone,
                v4_cidr_blocks=['10.0.0.0/24'],
        )

    def _restore_from_deployment(self, stack):
        resources = stack.export_stack().deployment['resources']
        for resource in resources:
            if resource['type'] != 'yandex:index/computeInstance:ComputeInstance':
                continue
            if resource['outputs']['status'] != 'running':
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
