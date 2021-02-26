'''
Scaling algorithm for GitLab runners fleet
'''


from collections.abc import Sequence

from .instance import InstanceParams


def require_instances(current: Sequence[InstanceParams]) -> Sequence[InstanceParams]:
    '''Yield instance params for all runners deemed necessary'''


def current_instances() -> Sequence[InstanceParams]:
    '''Yield instance params for all instances existing after previous run'''


def is_busy(params: InstanceParams):
    '''Check that runner instance is busy'''
