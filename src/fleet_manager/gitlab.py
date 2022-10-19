'''
Communicate with GitLab API
'''


import os
from itertools import chain
from time import sleep
from cirrus_run import CirrusAPI as GraphqlAPI

from .logging import log

try:
    from functools import cache
except ImportError:
    from functools import lru_cache
    cache = lru_cache(maxsize=None)


class GitLabAPI(GraphqlAPI):
    '''Reuse API object from cirrus-run project'''

    USER_AGENT = 'CI runners fleet manager <https://github.com/sio/gitlab-runners-fleet>'

    def __init__(self, api_token, runner_token, graphql_url, rest_url):
        self.runner_token = runner_token
        self.rest_url = rest_url
        super().__init__(url=graphql_url, token=api_token)

    def delete(self, *a, **ka):
        return self._requests.delete(*a, **ka)

    def post(self, *a, **ka):
        return self._requests.post(*a, **ka)

    @property
    @cache
    def namespace(api) -> str:
        '''Return GitLab namespace for CI runners'''
        username_query = '{currentUser {username}}'
        reply = api(username_query)
        try:
            username = reply['currentUser']['username']
        except (KeyError, TypeError):
            raise RuntimeError(f"GitLab couldn't match provided API token to any user account")
        return username

    def get_pending_jobs(api) -> int:
        '''Return number of GitLab CI jobs currently waiting for a runner'''
        jobs_query = '''
            query GetPendingJobs($namespace: ID!){
              namespace(fullPath: $namespace) {
                projects {
                  nodes {
                    running: pipelines(status: RUNNING) {
                      ...JobDetails
                    }
                    waiting: pipelines(status: WAITING_FOR_RESOURCE) {
                      ...JobDetails
                    }
                    pending: pipelines(status: PENDING) {
                      ...JobDetails
                    }
                  }
                }
              }
            }

            fragment JobDetails on PipelineConnection {
              nodes {
                jobs {
                  nodes {
                    detailedStatus {
                      text
                    }
                  }
                }
              }
            }
        '''
        pending_jobs = 0
        for project in api(
                jobs_query,
                params=dict(namespace=api.namespace)
        )['namespace']['projects']['nodes']:
            for pipeline in chain.from_iterable(status['nodes'] for status in project.values()):
                for job in pipeline['jobs']['nodes']:
                    if job['detailedStatus']['text'] == 'pending':
                        pending_jobs += 1
        return pending_jobs

    def get_project_ids(api):
        '''List project IDs associated with current namespace'''
        projects_query = '''
            query GetProjectIDs($namespace: ID!) {
                namespace(fullPath: $namespace) {
                    projects {
                        nodes{
                            id
                        }
                    }
                }
            }
        '''
        for project in api(
                projects_query,
                params=dict(namespace=api.namespace)
        )['namespace']['projects']['nodes']:
            gid = project['id']
            iid = int(gid.split('/')[-1])
            yield iid

    def get_runners(api, tags=None):
        params = dict()
        if tags:
            params['tag_list'] = tags
        url = f'{api.rest_url}/runners'
        response = api.get(url, params=params)
        response.raise_for_status()
        for runner in response.json():
            yield runner

    def unregister_runner(api, runner_id):
        runner = api.get(f'{api.rest_url}/runners/{runner_id}').json()
        for project in runner.get('projects', []):
            response = api.delete(f'{api.rest_url}/projects/{project["id"]}/runners/{runner_id}')
        response = api.delete(f'{api.rest_url}/runners/{runner_id}')
        if not response.ok:
            log.error('Could not unregister runner: %s', response.json())
        log.info(f'Unregistered runner %s #%s (was %s)', runner['description'], runner['id'], runner['status'])

    def assign_runner(api, runner_id, project_ids):
        response = api.get(f'{api.rest_url}/runners/{runner_id}')
        response.raise_for_status()

        assigned = set(project['id'] for project in response.json().get('projects', []))
        missing = set(project_ids).difference(assigned)
        if missing:
            log.info('Assigning runner #%s to projects %s', runner_id, missing)
        for project_id in missing:
            response = api.post(
                f'{api.rest_url}/projects/{project_id}/runners',
                data=dict(runner_id=runner_id)
            )
            response.raise_for_status()
        return bool(missing)

    def update_runner_assignments(api):
        projects = list(api.get_project_ids())
        updated = False
        for runner in api.get_runners(tags='private-runner'):
            if runner['status'] == 'online':
                if api.assign_runner(runner['id'], projects):
                    updated = True
            else:
                api.unregister_runner(runner['id'])
        if updated:
            sleep(10)  # allow newly assigned runners to pick up pending jobs
