'''
Pre-destroy cleanup of CI runner instances
'''


import subprocess
from data import InstanceParams


def cleanup(instance: InstanceParams, identity_file: str):
    '''
    Unregister GitLab runners on instances scheduled for destroying
    '''
    ssh = [
        'ssh',
        '-i', identity_file,
        '-o', 'UserKnownHostsFile=/dev/null',
        '-o', 'StrictHostKeyChecking=no',
        instance.ssh,
    ]
    ssh.extend(instance.cleanup)
    subprocess.run(ssh, check=True)
