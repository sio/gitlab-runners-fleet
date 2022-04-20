'''
Command line interface
'''

import argparse
from importlib import import_module
from pathlib import Path

from . import logging
from .config import Configuration
from .logging import log
from .scaling import ScalingConfig


def main(*a, **ka):
    args = parse_args(*a, **ka)
    config = Configuration(args.config)
    logging.setup(config.main.log_level)
    cloud = cloud_provider(config)
    print(cloud)



def cloud_provider(config):
    '''Initialize cloud provider object based on provided configuration'''
    class_path = config.main.cloud
    module_path, class_name = class_path.split(':', 1)
    if not module_path or not class_name:
        raise ValueError(f'invalid cloud provider class path: {class_path}')
    module = import_module(module_path)
    cloud_class = getattr(module, class_name)
    return cloud_class(
        scaling=ScalingConfig(**config.scaling),
        config=getattr(config, class_name, None),
    )


def dummy_cli():
    from pulumi import automation as auto
    stack = auto.create_or_select_stack(
        stack_name = 'dev',
        project_name = 'gitlab_runners_rewrite',
        program = cloud.pulumi,
    )
    stack.workspace.install_plugin(*cloud.PLUGIN)
    refresh = stack.refresh() # on_output=print)
    if sys.argv[1] == 'up':
        up = stack.up(on_output=print)
    elif sys.argv[1] == 'destroy':
        stack.destroy(on_output=print)


def parse_args(*a, **ka):
    parser = argparse.ArgumentParser(
        description='Maintain autoscaling fleet of GitLab CI runners',
        epilog='Licensed under Apache License, version 2.0'
    )
    parser.add_argument(
        'action',
        metavar='ACTION',
        choices=('up', 'destroy'),
        type=str.lower,
        help='infrastructure action to perform',
    )
    parser.add_argument(
        '--config',
        metavar='CONFIG',
        default='config.toml',
        nargs='?',
        type=Path,
        help='path to configuration file',
    )
    args = parser.parse_args(*a, **ka)
    if not args.config.exists():
        log.warning(f'Configuration file not found: {args.config}, continuing with defaults')
        args.config = None
    return args
