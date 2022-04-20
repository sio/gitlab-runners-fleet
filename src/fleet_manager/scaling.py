'''
Scaling configuration class with default values
'''


from dataclasses import dataclass
from typing import Optional


@dataclass
class ScalingConfig:
    jobs_per_instance: int = 2
    min_total_instances: int = 0
    max_total_instances: int = 10
    max_grow_instances: int = 2
    est_provisioning_minutes: int = 10
    min_billable_minutes: int = 10
    max_rerun_delay: int = 10
    max_idle_minutes: Optional[int] = None

    def __post_init__(self):
        if self.max_idle_minutes is None or self.max_idle_minutes < 0:
            self.max_idle_minutes = (self.min_billable_minutes
                                    -self.max_rerun_delay
                                    -self.est_provisioning_minutes)
