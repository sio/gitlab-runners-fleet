"""A Python Pulumi program"""

import json
import os
import pulumi

from data import InstanceParams
from destroy import cleanup
from instance import create, create_key
import scaling


config = pulumi.Config()
snapshot = config.get(os.environ['PULUMI_SNAPSHOT_OBJECT'])
if snapshot:
    params = json.loads(snapshot)
    if params:
        #cleanup(
        #    InstanceParams(**params),
        #    identity_file=os.environ['GITLAB_RUNNER_SSHKEY']
        #)
        pulumi.export('status_check', str(scaling.status(InstanceParams(**params))))


key = create_key()
server = create(InstanceParams('test-instance'), depends_on=[key,])


export = dict(
    name=server.name,
    cleanup=['/bin/touch', '/tmp/cleanup-worked'],
    metrics=server.ipv4_address.apply(lambda ip: f'http://{ip}:8080/metrics'),
    ssh=server.ipv4_address.apply(lambda ip: f'op@{ip}'),
)

pulumi.export(os.environ['PULUMI_SNAPSHOT_OBJECT'], export)
