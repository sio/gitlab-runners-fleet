"""A Python Pulumi program"""

import os
import pulumi

from data import InstanceParams
from destroy import cleanup
from instance import create, create_key


key = create_key()
server = create(InstanceParams('test-instance'), depends_on=[key,])


export = dict(
    name=server.name,
    cleanup=['/bin/touch', '/tmp/cleanup-worked'],
    metrics=server.ipv4_address.apply(lambda ip: f'http://{ip}:8080/metrics'),
    ssh=server.ipv4_address.apply(lambda ip: f'op@{ip}'),
)

pulumi.export(os.environ['PULUMI_SNAPSHOT_OBJECT'], export)
