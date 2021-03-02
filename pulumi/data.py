'''
Instance data holder
'''


from typing import Sequence
from dataclasses import dataclass, field
from enum import Enum, auto


@dataclass(frozen=True)
class InstanceParams:
    '''Class for keeping track of cloud server parameters'''
    name: str
    cleanup: Sequence[str] = field(default_factory=list)
    ssh: str = ''
    metrics: str = ''
    created_at: int = 0
    idle_since: int = 0


class InstanceStatus(Enum):
    BUSY = auto()
    READY = auto()
    IDLE = auto()
    PROVISIONING = auto()
    ERROR = auto()
