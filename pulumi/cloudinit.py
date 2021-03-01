'''
Render cloud-init userdata template
'''


import os
from pathlib import Path
from jinja2 import Template


def userdata(**fields):
    '''Render Jinja2 template'''
    userdata_path = os.environ['HCLOUD_USERDATA_TEMPLATE']
    with open(userdata_path) as f:
        template = Template(f.read())

    def instance_file(filename):
        '''Read file from the same directory as template'''
        directory = Path(userdata_path).parent
        with open(directory/filename) as f:
            return f.read()

    env = dict(instance_file=instance_file)
    env.update(fields)
    return template.render(env)
