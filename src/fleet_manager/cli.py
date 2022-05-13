'''
Command line interface
'''

import argparse
from importlib import import_module
from pathlib import Path

from pulumi import automation as auto

from . import logging
from .config import Configuration
from .gitlab import GitLabAPI
from .logging import log
from .scaling import ScalingConfig


def main(*a, **ka):
    args = parse_args(*a, **ka)

    config = Configuration(args.config)
    if args.action == 'down':
        args.action = 'up'
        config._data['scaling'].update(dict(
            min_total_instances = 0,
            max_total_instances = 0
        ))
    logging.setup(config.main.log_level)
    cloud = cloud_provider(config)
    log.debug(f'Initialized cloud provider object:\n{cloud}')

    stack = auto.create_or_select_stack(
        stack_name = config.pulumi.stack,
        project_name = config.pulumi.project,
        program = cloud.pulumi,
    )
    for plugin in config.pulumi.plugins:
        stack.workspace.install_plugin(*plugin)
    refresh = stack.refresh() # on_output=print)
    cloud.restore(stack)
    action = getattr(stack, args.action)
    action(on_output=print)


def cloud_provider(config):
    '''Initialize cloud provider object based on provided configuration'''
    class_path = config.main.cloud
    module_path, class_name = class_path.split(':', 1)
    if not module_path or not class_name:
        raise ValueError(f'invalid cloud provider class path: {class_path}')
    module = import_module(module_path)
    cloud_class = getattr(module, class_name)
    return cloud_class(
        config=getattr(config, class_name, None),
        gitlab=GitLabAPI(**config.gitlab),
        scaling=ScalingConfig(**config.scaling),
    )


def parse_args(*a, **ka):
    parser = argparse.ArgumentParser(
        description='Maintain autoscaling fleet of GitLab CI runners',
        epilog='Licensed under Apache License, version 2.0'
    )
    parser.add_argument(
        'action',
        metavar='ACTION',
        choices=('up', 'destroy', 'down'),
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
