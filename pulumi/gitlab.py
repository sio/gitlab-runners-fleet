'''
Communicate with GitLab API
'''


from cirrus_run import CirrusAPI as GraphqlAPI


class GitLabAPI(GraphqlAPI):
    '''Reuse API object from cirrus-run project'''
    DEFAULT_URL = 'https://gitlab.com/api/graphql'
    USER_AGENT = 'CI runners fleet manager <https://github.com/sio/gitlab-runners-fleet>'


def get_pending_jobs() -> int:
    '''Return number of GitLab CI jobs currently waiting for a runner'''
