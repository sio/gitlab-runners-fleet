'''
Scaling algorithm for GitLab runners fleet
'''

import coolname
import json
import math
import os
import pulumi
import requests
from dataclasses import replace
from typing import Sequence

import gitlab
import timestamp
from data import InstanceParams, InstanceStatus


JOBS_PER_INSTANCE = 2  # This is used to calculate the number of required instances.
                       # The value should be less than or equal to the maximum number of
                       # concurrent jobs allowed per instance

MAX_TOTAL_INSTANCES = 10  # How many instances may exist at once
MAX_GROW_INSTANCES = 2    # How many instances may be added at once

# Idle timeout length calculations
EST_PROVISIONING_MINUTES = 10
MIN_BILLABLE_MINUTES = 60
MAX_RERUN_DELAY = 10
MAX_IDLE_MINUTES = MIN_BILLABLE_MINUTES - MAX_RERUN_DELAY - EST_PROVISIONING_MINUTES


def get_status(instance: InstanceParams) -> InstanceStatus:
    '''Detect instance status'''
    now = timestamp.now()
    try:
        response = requests.get(instance.metrics)
        response.raise_for_status()
        metrics = response.json()
    except Exception:
        if now - instance.created_at < EST_PROVISIONING_MINUTES * 60:
            return InstanceStatus.PROVISIONING
        return InstanceStatus.ERROR
    if metrics.get('gitlab_runner_jobs', 0) > 0:
        return InstanceStatus.BUSY
    if instance.idle_since and now - instance.idle_since > MAX_IDLE_MINUTES * 60:
        return InstanceStatus.IDLE
    return InstanceStatus.READY


def calculate_actions():
    '''Calculate infra scaling actions'''
    actions = calculate_actions_keep_delete()
    actions.update({
        'CREATE': {
            InstanceStatus.NOT_EXISTS: set(),
        },
    })

    jobs_capacity = JOBS_PER_INSTANCE * (
        len(actions['KEEP'][InstanceStatus.PROVISIONING]) +
        len(actions['KEEP'][InstanceStatus.READY])
    )
    jobs_required = gitlab.get_pending_jobs()
    instances_required = max(0, int(math.ceil((jobs_required - jobs_capacity)/JOBS_PER_INSTANCE)))

    if not instances_required:
        return actions

    names_taken = set()
    count_instances = 0
    for action, status_set in actions.items():
        for instances in status_set.values():
            for instance in instances:
                names_taken.add(instance.name)
                if action != 'DELETE':
                    count_instances += 1

    instances_to_add = min(
        instances_required,
        MAX_GROW_INSTANCES,
        MAX_TOTAL_INSTANCES - count_instances,
    )

    for _ in range(instances_to_add):
        new_name = ''
        while not new_name or new_name in names_taken:
            new_name = 'ci-' + coolname.generate_slug(2)
        actions['CREATE'][InstanceStatus.NOT_EXISTS].add(
            InstanceParams(name=new_name, created_at=timestamp.now())
        )
    return actions


def calculate_actions_keep_delete():
    '''Decide which of previously existing instances to keep and which to delete'''
    actions = {
        'KEEP': {
            InstanceStatus.PROVISIONING: set(),
            InstanceStatus.READY: set(),
            InstanceStatus.BUSY: set(),
        },
        'DELETE': {
            InstanceStatus.ERROR: set(),
            InstanceStatus.IDLE: set(),
        },
    }
    config = pulumi.Config()
    snapshot = config.get(os.environ['PULUMI_SNAPSHOT_OBJECT'])
    if not snapshot:
        return actions
    previous_state = json.loads(snapshot)
    if not previous_state:
        return actions

    for params in previous_state:
        instance = InstanceParams(**params)
        status = get_status(instance)
        if status == InstanceStatus.READY and not instance.idle_since:
            instance = replace(instance, idle_since=timestamp.now())
        if status == InstanceStatus.BUSY and instance.idle_since:
            instance = replace(instance, idle_since=0)
        if status in actions['KEEP']:
            actions['KEEP'][status].add(instance)
        else: # InstanceStatus.ERROR, InstanceStatus.IDLE
            actions['DELETE'][status].add(instance)
    return actions
