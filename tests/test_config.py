import pytest
from fleet_manager.config import AttributeDictReader, Configuration


def test_attribute_access():
    tree = AttributeDictReader(dict(
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


def test_reading_default_config():
    config = Configuration()
    assert config.scaling.jobs_per_instance == 2
    assert config.YandexCloud.preemptible_instances == True
    assert config.pulumi.stack == 'stackname'


def test_configuration_merging():
    config = Configuration('tests/config_override.toml')
    assert config.scaling.jobs_per_instance == 2
    assert config.YandexCloud.preemptible_instances == True
    assert config.pulumi.stack == 'stackname'
    assert config.YandexCloud.memory_gb == 8
    assert config.main.cloud == 'NonExistentProvider'
