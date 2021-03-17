'''Auto scaling fleet of GitLab CI runners'''

import os
import pulumi
from itertools import chain

import scaling
from destroy import cleanup
from instance import create, create_key


pulumi.log.debug('Create SSH key in cloud account')
key = create_key()


pulumi.log.debug('Calculate actions for existing instances')
actions = scaling.calculate_actions()


pulumi.log.debug('Execute cleanup actions on machines scheduled for deletion')
for status, instances in actions['DELETE'].items():
    for instance in instances:
        pulumi.log.info(f'Deleting instance {instance.name}: {status}')
        try:
            cleanup(instance, identity_file=os.environ['GITLAB_RUNNER_SSHKEY'])
        except Exception:
            pass


pulumi.log.debug('Create servers')
export = []
for status, instances in chain(actions['KEEP'].items(), actions['CREATE'].items()):
    for instance in instances:
        pulumi.log.debug(f'Create server: {instance.name}')
        server = create(instance, depends_on=[key,])
        export.append(dict(
            name=instance.name,
            cleanup=instance.cleanup or [
                'sudo', '-u', 'gitlab-runner', '/etc/gitlab-runner-custom/unregister.sh'
            ],
            ssh=server.ipv4_address.apply(lambda ip: f'op@{ip}'),
            metrics=server.ipv4_address.apply(lambda ip: f'http://{ip}:8080/metrics'),
            created_at=instance.created_at,
            idle_since=instance.idle_since,
        ))


pulumi.log.debug('Export infrastructure snapshot')
pulumi.export(
    os.environ['PULUMI_SNAPSHOT_OBJECT'],
    sorted(
        export,
        key=lambda i: (i['created_at'], i['name']),
    )
)
