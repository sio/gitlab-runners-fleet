'''
Scaling algorithm for GitLab runners fleet
'''


import requests
from datetime import datetime
from typing import Sequence

from data import InstanceParams, InstanceStatus


EST_PROVISIONING_MINUTES = 10
MIN_BILLABLE_MINUTES = 60
MAX_RERUN_DELAY = 10
MAX_IDLE_MINUTES = MIN_BILLABLE_MINUTES - MAX_RERUN_DELAY - EST_PROVISIONING_MINUTES


def require_instances(current: Sequence[InstanceParams]) -> Sequence[InstanceParams]:
    '''Yield instance params for all runners deemed necessary'''


def current_instances() -> Sequence[InstanceParams]:
    '''Yield instance params for all instances existing after previous run'''


def status(instance: InstanceParams) -> InstanceStatus:
    '''Detect instance status'''
    now = datetime.utcnow().timestamp()
    if now - instance.created_at < EST_PROVISIONING_MINUTES * 60:
        return InstanceStatus.PROVISIONING
    try:
        response = requests.get(instance.metrics)
        response.raise_for_status()
        metrics = response.json()
    except Exception:
        return InstanceStatus.ERROR
    if metrics.get('gitlab_runner_jobs', 0) > 0:
        return InstanceStatus.BUSY
    if instance.idle_since and now - instance.idle_since > MAX_IDLE_MINUTES * 60:
        return InstanceStatus.IDLE
    return InstanceStatus.READY
