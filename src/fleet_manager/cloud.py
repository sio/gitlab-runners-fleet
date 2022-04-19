'''
Abstract base classes for uniform interaction with different cloud providers
'''


import math
from abc import ABC, abstractmethod
from dataclasses import dataclass
from enum import Enum, auto

import coolname

from . import gitlab, timestamp

from __future__ import annotations  # https://stackoverflow.com/a/52699243
from typing import Optional


class InstanceStatus(Enum):
    NEW = auto()            # not deployed yet
    PROVISIONING = auto()   # deployed but is not ready to accept jobs yet
    READY = auto()          # ready to accept jobs but is not currently executing any
    BUSY = auto()           # currently executing one or more jobs
    IDLE = auto()           # has not been used for a long time, or has reached max allowed age
    DESTROYING = auto()     # cleanup was completed, instance is ready to be destroyed
    ERROR = auto()          # irrecoverable error, instance needs to be destroyed
status = InstanceStatus


@dataclass
class CloudInstance(ABC):
    '''Abstract class for a cloud compute instance'''

    cloud: CloudProvider
    name: str
    status: InstanceStatus = InstanceStatus.NEW
    idle_since: int = 0
    created_at: int = 0

    @abstractmethod
    def update_status(self):
        '''Write updated values to self.status, self.idle_since'''
        if self.status == status.READY and not self.idle_since:
            self.idle_since = timestamp.now()
        if self.status == status.BUSY and self.idle_since:
            self.idle_since = 0

    @abstractmethod
    def create(self):
        '''Create cloud instance corresponding to this object'''
        self.created_at = timestamp.now()
        self.status = status.PROVISIONING

    @abstractmethod
    def cleanup(self):
        '''
        Prepare instance for deletion:
            - Unregister GitLab runners
            - Remove cloud firewall rules
            - etc.

        After successful cleanup this method must set self.status to DESTROYING
        '''
        self.status = status.DESTROYING

    def __repr__(self):
        return f'<{self.__class__.__name__}: {self.name} ({self.status.name})>'


class CloudProvider(ABC):
    '''Abstract class for cloud provider'''

    _instance_cls: CloudInstance
    _namelog_maxlen_multiplier = 50

    def __init__(self, scaling_config=None):
        self.instances: set[CloudInstance] = set()
        self.scaling = ScalingConfig()
        if scaling_config is not None:
            self.scaling = scaling_config
        self._names_seen = set()

    def new(instance_name: Optional[str] = None) -> CloudInstance:
        '''Return new cloud instance object'''
        if instance_name is None:
            instance_name = self._new_name()
        self._check_name(instance_name)
        self._names_seen.add(instance_name)

        instance = self._instance_cls(cloud=self, name=name)
        self.instances.add(instance)
        return instance

    def _check_name(self, instance_name):
        '''Validate name for new instance'''
        if len(self._names_seen) < len(self.instances) \
        or len(self._names_seen) > len(self.instances) * self._namelog_maxlen_multiplier:
            self._names_seen = set(i.name for i in self.instances)
        if instance_name in self._names_seen:
            raise ValueError(f'instance name is in use or was recently in use: {instance_name}')

    def _new_name(self, prefix=''):
        '''Generate new instance name'''
        name = None
        while not name or name in self._names_seen:
            name = prefix + coolname.generate_slug(2)
        return name

    def pulumi(self):
        '''Inline program for Pulumi Automation API'''
        if not self.instances:
            return
        self.setup()
        for instance in sorted(self.instances):
            instance.create()

    def scale(self):
        '''Calculate scaling actions for cloud instances'''
        for instance in self.instances:
            instance.update_status()
            if instance.status in {
                    status.ERROR,
                    status.IDLE,
            }:
                instance.cleanup()
            if instance.status in {
                    status.DESTROYING,
                    status.ERROR,
            }:
                # pulumi will destroy everything it wasn't explicitly asked to keep
                self.instances.remove(instance)

        scaling = self.scaling
        jobs_pending = gitlab.get_pending_jobs()
        jobs_capacity = scaling.jobs_per_instance * len(
                i for i in self.instances if i.status in {status.PROVISIONING, status.READY}
            )
        instances_required = max(
                0,
                int(math.ceil(
                    (jobs_pending - jobs_capacity) / scaling.jobs_per_instance
                ))
            )
        instances_to_add = min(
                instances_required,
                scaling.max_grow_instances,
                scaling.max_total_instances - len(self.instances),
            )
        for _ in range(instances_to_add):
            self.new()

    @abstractmethod
    def setup():
        '''
        Ensure that cloud provider is ready for creating instances:
            - Create required SSH keys
            - Configure cloud networking
            - Configure cloud NAT/firewall
            - etc.
        '''
