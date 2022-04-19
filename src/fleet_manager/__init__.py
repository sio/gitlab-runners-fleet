import sys

def dummy_cli():
    print(sys.argv)

    from . import logging
    logging.setup()

    from .scaling import ScalingConfig
    scaling_config = ScalingConfig(
        min_total_instances = 1
    )

    from .yandex import YandexCloud
    cloud = YandexCloud(scaling_config)

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
