"""A Python Pulumi program"""

import json
import os
import pulumi

from data import InstanceParams
from destroy import cleanup
from instance import create, create_key


config = pulumi.Config()
snapshot = config.get(os.environ['PULUMI_SNAPSHOT_OBJECT'])
if snapshot:
    params = json.loads(snapshot)
    if params:
        cleanup(
            InstanceParams(**params),
            identity_file=os.environ['RUNNER_SSH_KEY']
        )


create_key()
server = create(InstanceParams('test-instance'))  # TODO: add dependency on key object?


export = dict(
    name=server.name,
    endpoint='',
    ssh=server.ipv4_address.apply(lambda ip: f'root@{ip}'),
    cleanup=['/bin/touch', '/tmp/cleanup-worked'],
)

pulumi.export(os.environ['PULUMI_SNAPSHOT_OBJECT'], export)
