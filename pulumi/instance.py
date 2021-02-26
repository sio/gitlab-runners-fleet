'''
Instance definition for a single GitLab runner
'''


from dataclasses import dataclass


@dataclass(frozen=True)
class InstanceParams:
    '''Class for keeping track of cloud server parameters'''
    name: str
    endpoint: str = ''


def create(params: InstanceParams):
    '''Create GitLab runner cloud instance'''
