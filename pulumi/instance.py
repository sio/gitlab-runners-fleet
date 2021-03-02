'''
Instance definition for a single GitLab CI runner
'''


import os
from dataclasses import dataclass

from pulumi import ResourceOptions
import pulumi_hcloud as hcloud

import cloudinit
from data import InstanceParams
from scaling import (
    require_instances,
)


SSH_KEY_NAME = 'Pulumi SSH key for gitlab-runners-fleet'


def create_key():
    '''Upload SSH key to Hetzner account'''
    key = hcloud.SshKey(
        resource_name='pulumi-key',
        public_key=read_file(os.environ['GITLAB_RUNNER_SSHKEY'] + '.pub'),
        name=SSH_KEY_NAME,
    )
    return key


def create(params: InstanceParams, depends_on=None):
    '''Create GitLab CI runner cloud instance'''
    if not depends_on:
        depends_on = []
    server = hcloud.Server(
        resource_name=params.name,
        name=params.name,
        server_type='cx11',
        image='debian-10',
        ssh_keys=[SSH_KEY_NAME,],
        user_data=cloudinit.userdata(
            pubkey=read_file(os.environ['GITLAB_RUNNER_SSHKEY'] + '.pub'),
            gitlab_runner_token=os.environ['GITLAB_RUNNER_TOKEN'],
        ),
        opts=ResourceOptions(depends_on=depends_on),
    )
    return server


def read_file(path):
    '''Read contents of a small file'''
    with open(path) as f:
        return f.read()
