'''
Instance definition for a single GitLab CI runner
'''


from dataclasses import dataclass

import pulumi_hcloud as hcloud

from data import InstanceParams
from scaling import (
    require_instances,
)


def create(params: InstanceParams):
    '''Create GitLab CI runner cloud instance'''
    server = hcloud.Server(
        resource_name=params.name,
        server_type='cx11',
        image='debian-10'
    )
    return server
