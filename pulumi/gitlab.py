'''
Communicate with GitLab API
'''


def get_pending_jobs() -> int:
    '''Return number of GitLab CI jobs currently waiting for a runner'''
