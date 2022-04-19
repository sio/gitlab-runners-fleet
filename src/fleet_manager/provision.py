'''
Working with instance configuration
'''

import atexit
from pathlib import Path
from pkg_resources import (
        cleanup_resources,
        resource_exists,
        resource_filename,
        resource_isdir,
)

try:
    from functools import cache
except ImportError:
    from functools import lru_cache
    cache = lru_cache(maxsize=None)

import jinja2

from .logging import log


def resolve_path(path):
    '''
    Resolve path that might be any of the following:
        - an absolute path to file on filesystem
        - a relative path to file on filesystem
        - a relative path to package resource
    '''
    if resource_exists(__package__, path):
        dirname, filename = path.rsplit('/', 1)  # resource paths must not be manipulated via os.path
        if resource_exists(__package__, dirname) and resource_isdir(__package__, dirname):
            path = Path(resource_filename(__package__, dirname)) / filename
        else:
            path = Path(resource_filename(__package__, path))
        atexit.register(cleanup_resources)
    else:
        path = Path(path)
    return path


def template(path):
    '''Return Jinja2 template corresponding to provided path'''
    path = resolve_path(path)
    if not path.exists:
        raise ValueError(f'template not found in path: {path}')
    j2 = jinja_environment(path.parent)
    return j2.get_template(path.name)


@jinja2.pass_environment
def static_file(j2, name):
    '''Jinja2 function to render contents of a static file from template directory'''
    directory = Path(j2.globals['template_directory'])
    filepath = directory / name
    with filepath.open() as f:
        return f.read()


@cache
def jinja_environment(template_directory):
    '''Create Jinja2 environment for the template directory'''
    log.debug('Creating new Jinja2 environment for %s', template_directory)
    env = jinja2.Environment(loader=jinja2.FileSystemLoader(template_directory))
    env.globals.update(dict(
        static_file = static_file,
        template_directory = str(template_directory),
    ))
    return env
