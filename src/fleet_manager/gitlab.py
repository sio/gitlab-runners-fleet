'''
Communicate with GitLab API
'''


import os
from itertools import chain
from cirrus_run import CirrusAPI as GraphqlAPI

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

    def get_project_ids():
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
