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
            name = self.name,
            scheduling_policy = yandex.ComputeInstanceSchedulingPolicyArgs(
                preemptible = True,
            ),
            metadata = {
                'user-data': template('provisioning/cloudinit.yml.j2').render(
                    pubkey = 'publickey', # TODO
                    gitlab_runner_token = 'gitlab_runner_token', # TODO
                ),
            },
            resources = yandex.ComputeInstanceResourcesArgs(
                cores = 2,
                memory = 4,
            ),
            zone = 'ru-central1-a',  # TODO: allow configuration
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
