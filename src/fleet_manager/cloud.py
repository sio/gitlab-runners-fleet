'''
Abstract base classes for uniform interaction with different cloud providers
'''


from __future__ import annotations  # https://stackoverflow.com/a/52699243
from typing import Optional

import math
from abc import ABC, abstractmethod
from dataclasses import dataclass, asdict, fields
from enum import Enum, auto

import coolname
import pulumi

from . import gitlab, timestamp
from .logging import log
from .scaling import ScalingConfig


class InstanceStatus(Enum):
    NEW = auto()            # not deployed yet
    PROVISIONING = auto()   # deployed but is not ready to accept jobs yet
    READY = auto()          # ready to accept jobs but is not currently executing any
    BUSY = auto()           # currently executing one or more jobs
    IDLE = auto()           # has not been used for a long time, or has reached max allowed age
    DESTROYING = auto()     # cleanup was completed, instance is ready to be destroyed
    ERROR = auto()          # irrecoverable error, instance needs to be destroyed
status = InstanceStatus


@dataclass(eq=False)
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

    def __init__(self, config=None, gitlab=None, scaling=None):
        self.config = config or dict()
        self.gitlab = gitlab or dict()
        self.scaling = scaling or ScalingConfig()
        self.instances: set[CloudInstance] = set()
        self._names_seen = set()

    def __repr__(self):
        return f'<{self.__class__.__name__} scaling={asdict(self.scaling)} config={dict(self.config)}>'

    def new(self, instance_name: Optional[str] = None) -> CloudInstance:
        '''Return new cloud instance object'''
        if instance_name is None:
            instance_name = self._new_name()
        self._check_name(instance_name)
        self._names_seen.add(instance_name)

        instance = self._instance_cls(cloud=self, name=instance_name)
        instance.created_at = timestamp.now()
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
        self.scale()
        if self.instances:
            self.setup()
        for instance in self.instances:
            instance.create()
        self.save()

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
        jobs_capacity = scaling.jobs_per_instance * len([
                None for i in self.instances if i.status in {status.PROVISIONING, status.READY}
            ])
        instances_required = max(
                0,
                int(math.ceil(
                    (jobs_pending - jobs_capacity) / scaling.jobs_per_instance
                ))
            )
        instances_to_add = max(
            scaling.min_total_instances - len(self.instances),
            min(
                instances_required,
                scaling.max_grow_instances,
                scaling.max_total_instances - len(self.instances),
            ),
        )
        for _ in range(instances_to_add):
            self.new()

    @abstractmethod
    def setup(self):
        '''
        Ensure that cloud provider is ready for creating instances:
            - Create required SSH keys
            - Configure cloud networking
            - Configure cloud NAT/firewall
            - etc.
        '''

    def save(self):
        '''
        Save all instances to persistent storage (Pulumi stack)
        '''
        export = {}
        for instance in self.instances:
            data = {}
            for field in fields(instance):
                if field.name in {'cloud', 'name', 'status'}:
                    continue
                data[field.name] = getattr(instance, field.name)
            export[instance.name] = data
        log.debug('Saving stack output: %s', export)
        pulumi.export(self.__class__.__name__, export)

    def restore(self, stack):
        '''
        Restore instance list from persistent storage (Pulumi stack)
        '''
        if self.instances:
            raise RuntimeError('can not restore over existing instances')
        self._restore_from_stack_output(stack)
        if not self.instances:  # stack output may be lost on crash
            self._restore_from_deployment(stack)

    def _restore_from_deployment(self, stack):
        '''Restore instance list from current deployment'''
        log.warning(
            'Restoring instance list from deployment is not implemented for %s',
            self.__class__.__name__,
        )

    def _restore_from_stack_output(self, stack):
        '''Restore instance list from previous stack output'''
        export = stack.outputs().get(self.__class__.__name__)
        if export is None:
            log.debug('No saved instances in stack output for %s', self.__class__.__name__)
            return
        log.debug('Restoring instances from stack output: %s', export)
        for name, params in export.value.items():
            instance = self._instance_cls(name=name, cloud=self, **params)
            self.instances.add(instance)
