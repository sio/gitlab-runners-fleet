'''
Assign created runners to all projects in current namespace
'''


from time import sleep
from gitlab import api, get_project_ids


GITLAB_REST_API = 'https://gitlab.com/api/v4'


def get_runner_ids():
    '''Use REST API to yield IDs of CI runners'''
    url = GITLAB_REST_API + '/runners'
    response = api.get(url, params=dict(tag_list='private-runner', status='online'))
    response.raise_for_status()
    for runner in response.json():
        yield runner['id']


def assign_runner(runner_id, project_ids):
    '''
    Assign runner to multiple GitLab projects
    Return boolean to indicate if any changes were made
    '''
    url = f'{GITLAB_REST_API}/runners/{runner_id}'
    response = api.get(url)
    response.raise_for_status()

    assigned_projects = set(project['id'] for project in response.json()['projects'])
    missing_projects = set(project_ids).difference(assigned_projects)
    if missing_projects:
        print(f'Assigning runner #{runner_id} to projects {missing_projects}')
    for project_id in missing_projects:
        url = f'{GITLAB_REST_API}/projects/{project_id}/runners'
        api.post(url, data=dict(runner_id=runner_id)).raise_for_status()  # requires API write access
    return bool(missing_projects)


if __name__ == '__main__':
    print('Updating runner assignments...')
    project_ids = list(get_project_ids())
    updated = False
    for runner_id in get_runner_ids():
        if assign_runner(runner_id, project_ids):
            updated = True
    if updated:
        sleep(10)  # allow newly assigned runners to pick up pending jobs
