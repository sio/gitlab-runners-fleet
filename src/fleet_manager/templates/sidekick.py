#!/usr/bin/env python3
'''
GitLab runner needs to be kicked sometimes to receive jobs from projects it was
assigned to after being started

Usecase:

- Project A has some CI jobs stuck because there are no runners available
- Runner assigned to another project B comes online
- Runner gets assigned to project A after starting up
- Runner never pick up the stuck job, and GitLab does not force it unless some
  new CI event happens server side
'''

import json
import urllib.request
import sys
from dataclasses import dataclass
from datetime import datetime, timezone
from time import sleep
from subprocess import run

def kick():
    print('Restarting idle gitlab-runner')
    run('systemctl try-restart gitlab-runner.service'.split(), stdout=sys.stdout, stderr=sys.stderr)

def is_idle() -> bool:
    try:
        response = urllib.request.urlopen('http://localhost:8080/metrics')
        if response.status != 200:
            raise ValueError(f'HTTP {response.status} for metrics wrapper endpoint')
        metrics = json.load(response)
        return metrics['gitlab_runner_jobs'] == 0 and metrics['gitlab_runner_jobs_total'] == 0
    except Exception as exc:
        print(exc)
    return True

def now():
    return datetime.now(tz=timezone.utc)

@dataclass
class GrowingDelay:
    minimum: int = 1
    maximum: int = 30
    step: int = 1
    value: int = 0

    def __call__(self):
        return self.increase()

    def __post_init__(self):
        if self.value < self.minimum:
            self.value = self.minimum

    def increase(self):
        if self.value < self.minimum:
            self.value = self.minimum
        elif self.value < self.maximum:
            self.value = min(self.maximum, self.value + self.step)
        return self.value

def main():
    start_at = now()
    idle_since = None
    idle_check_delay = GrowingDelay(minimum=1, maximum=30, step=1)
    restart_delay = GrowingDelay(minimum=15, maximum=60, step=5)

    while True:
        sleep(idle_check_delay())
        if not is_idle():
            print('Runner has proved to be working, exiting the sidekick script')
            sys.exit(0)
        if not idle_since:
            idle_since = now()
        idle_for = int((now() - idle_since).total_seconds())
        print(f'Runner has been idle for {idle_for} seconds (since {idle_since})')
        if idle_for > restart_delay.value:
            restart_delay.increase()
            kick()
            idle_since = now()

if __name__ == '__main__':
    main()
