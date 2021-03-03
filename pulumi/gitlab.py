'''
Communicate with GitLab API
'''


import os
from itertools import chain
from cirrus_run import CirrusAPI as GraphqlAPI


class GitLabAPI(GraphqlAPI):
    '''Reuse API object from cirrus-run project'''
    DEFAULT_URL = 'https://gitlab.com/api/graphql'
    USER_AGENT = 'CI runners fleet manager <https://github.com/sio/gitlab-runners-fleet>'


api = GitLabAPI(os.environ['GITLAB_API_TOKEN'])


def get_namespace() -> str:
    '''Return GitLab namespace for CI runners'''
    username_query = '{currentUser {username}}'
    username = api(username_query)['currentUser']['username']
    return username


def get_pending_jobs() -> int:
    '''Return number of GitLab CI jobs currently waiting for a runner'''
    jobs_query = '''
        query GetPendingJobs($namespace: ID!){
          namespace(fullPath: $namespace) {
            projects {
              nodes {
                waiting: pipelines(status: WAITING_FOR_RESOURCE) {  # TODO: Add RUNNING here?
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
    for project in api(jobs_query, params=dict(namespace=get_namespace()))['namespace']['projects']['nodes']:
        for pipeline in chain(project['waiting']['nodes'], project['pending']['nodes']):
            for job in pipeline['jobs']['nodes']:
                if job['detailedStatus']['text'] == 'pending':
                    pending_jobs += 1
    return pending_jobs
