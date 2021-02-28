'''
Instance data holder
'''


from typing import Sequence
from dataclasses import dataclass, field


@dataclass(frozen=True)
class InstanceParams:
    '''Class for keeping track of cloud server parameters'''
    name: str
    endpoint: str = ''
    ssh: str = ''
    cleanup: Sequence[str] = field(default_factory=list)
