'''
Configure logging
'''


import os
import logging
log = logging.getLogger(__package__)


def setup(level=None):
    '''
    Setup logging with basic defaults
    '''
    template = '%(levelname)-8s %(message)s'
    logging.basicConfig(format=template)
    if level is None:
        if os.getenv('DEBUG'):
            level = logging.DEBUG
        else:
            level = logging.WARNING
    log.level = level
    log.debug('Starting logging: %s', log)
