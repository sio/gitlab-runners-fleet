"""A Python Pulumi program"""

import os
import pulumi

from data import InstanceParams
from destroy import cleanup
from instance import create, create_key


create_key()
server = create(InstanceParams('test-instance'))  # TODO: add dependency on key object?
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

config = pulumi.Config()  # TODO: add Makefile step to transfer state from output to config (outside of Pulumi)
config.get('instance_params')  # TODO: something like 'pulumi stack output NAME --json|pulumi config set NAME $(cat /dev/stdin)'
