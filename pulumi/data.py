'''
Instance data holder
'''


from dataclasses import dataclass


@dataclass(frozen=True)
class InstanceParams:
    '''Class for keeping track of cloud server parameters'''
    name: str
    endpoint: str = ''
