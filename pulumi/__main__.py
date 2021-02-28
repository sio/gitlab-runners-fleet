"""A Python Pulumi program"""

import os
import pulumi

from data import InstanceParams
from destroy import cleanup
from instance import create, create_key


create_key()
server = create(InstanceParams('test-instance'))
server = create(InstanceParams('second-instance'))

#stack = pulumi.StackReference(pulumi.get_stack())
#cleanup(
#    InstanceParams(**stack.get_output('instance_params')),
#    identity_file=os.environ['RUNNER_SSH_KEY']
#)
export = dict(
    name=server.name,
    endpoint='',
    ssh=server.ipv4_address.apply(lambda ip: f'root@{ip}'),
    cleanup=['/bin/touch', '/tmp/cleanup-worked'],
)

pulumi.export('instance_params', export)
