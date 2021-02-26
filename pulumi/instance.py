'''
Instance definition for a single GitLab runner
'''


from dataclasses import dataclass


from data import InstanceParams


def create(params: InstanceParams):
    '''Create GitLab runner cloud instance'''
