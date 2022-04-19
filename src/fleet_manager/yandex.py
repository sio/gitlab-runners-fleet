'''
Yandex Cloud
'''

import pulumi_yandex as yandex

from .cloud import CloudProvider, CloudInstance, status
from .provision import template


class YandexInstance(CloudInstance):
    '''Yandex Cloud VPS'''

    def update_status(self):
        '''Write updated values to self.status, self.idle_since'''
        super().update_status()

    def create(self):
        '''Create cloud instance corresponding to this object'''
        yandex.ComputeInstance(
            resource_name = self.name,
            hostname = self.name,
            scheduling_policy = yandex.ComputeInstanceSchedulingPolicyArgs(
                preemptible = True,
            ),
            metadata = {
                'user-data': template('provisioning/cloudinit.yml.j2').render(
                    pubkey = 'publickey', # TODO
                    gitlab_runner_token = 'gitlab_runner_token', # TODO
                ),
            },
            zone = 'ru-central1-a',  # TODO: allow configuration
            resources = yandex.ComputeInstanceResourcesArgs(
                cores = 2,  # TODO: allow configuration
                memory = 4,
            ),
            boot_disk=yandex.ComputeInstanceBootDiskArgs(
                auto_delete=True,
                initialize_params=yandex.ComputeInstanceBootDiskInitializeParamsArgs(
                    size=15,  # TODO: allow configuration
                    image_id=yandex.get_compute_image(family='debian-11').image_id,  # TODO: allow configuration
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

    PLUGIN = ('yandex', 'v0.13.0')
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
                zone='ru-central1-a', # TODO: allow configuration
                v4_cidr_blocks=['10.0.0.0/24'],
        )
