import os
import pytest
from fleet_manager.config import ConfigurationTree, Configuration


def test_attribute_access():
    tree = ConfigurationTree(dict(
        hello = 'world',
        foo = 'bar',
        alice = dict(
            age = 30,
            location = 'NZ',
        )
    ))
    assert tree.hello == 'world'
    assert tree.foo == 'bar'
    assert tree.alice.location == 'NZ'
    assert tree['hello'] == 'world'
    assert tree['foo'] == 'bar'
    assert tree['alice']['location'] == 'NZ'


def test_reading_default_config():
    config = Configuration()
    assert config.scaling.jobs_per_instance == 3
    assert config.YandexCloud.preemptible_instances == True
    assert config.pulumi.stack == 'stackname'
    assert config['pulumi']['stack'] == 'stackname'


def test_configuration_merging():
    config = Configuration('tests/config_override.toml')
    assert config.scaling.jobs_per_instance == 3
    assert config.YandexCloud.preemptible_instances == True
    assert config.pulumi.stack == 'stackname'
    assert config.YandexCloud.memory_gb == 8
    assert config.main.cloud == 'NonExistentProvider'
    assert config['pulumi']['stack'] == 'stackname'
    assert config['YandexCloud']['memory_gb'] == 8
    assert config['main']['cloud'] == 'NonExistentProvider'


def test_environment_values():
    env = os.environ.copy()
    config = Configuration('tests/config_override.toml')
    os.environ['TEST_VARIABLE']= 'hello world'
    assert config.main.environment_test == 'hello world'
    assert config['main']['environment_test'] == 'hello world'
    os.environ['TEST_VARIABLE']= 'CHANGED VALUE'
    assert config.main.environment_test == 'CHANGED VALUE'
    assert config['main']['environment_test'] == 'CHANGED VALUE'
    os.environ = env
