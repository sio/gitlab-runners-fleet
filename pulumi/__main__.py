'''Auto scaling fleet of GitLab CI runners'''

import os
import pulumi
from itertools import chain

import scaling
from destroy import cleanup
from instance import create, create_key


key = create_key()
actions = scaling.calculate_actions()

for status, instances in actions['DELETE'].items():
    for instance in instances:
        try:
            cleanup(instance, identity_file=os.environ['GITLAB_RUNNER_SSHKEY'])
        except Exception:
            pass

export = []
for status, instances in chain(actions['KEEP'].items(), actions['CREATE'].items()):
    for instance in instances:
        server = create(instance, depends_on=[key,])
        export.append(dict(
            name=instance.name,
            cleanup=('/bin/touch', '/tmp/cleanup-worked'),  # TODO: unregister runners
            ssh=server.ipv4_address.apply(lambda ip: f'op@{ip}'),
            metrics=server.ipv4_address.apply(lambda ip: f'http://{ip}:8080/metrics'),
            created_at=instance.created_at,
            idle_since=instance.idle_since,
        ))
pulumi.export(os.environ['PULUMI_SNAPSHOT_OBJECT'], export)
