'''
Configuration of fleet_manager app
'''


from collections.abc import Mapping

import toml
from pkg_resources import resource_string


class Configuration:
    '''TOML based application configuration'''

    def __init__(self, config_path=None, defaults='config_default.toml'):
        if config_path is None:
            custom = {}
        else:
            with open(config_path) as f:
                custom = toml.load(f)
        defaults = toml.loads(resource_string(__package__, defaults).decode())

        # only first-level dictionaries are properly merged,
        # deeper nested levels get overwritten
        merged = dict()
        for key in custom.keys():
            if  key in defaults \
            and isinstance(custom[key], Mapping) \
            and isinstance(defaults[key], Mapping):
                defaults[key].update(custom[key])
                merged[key] = defaults[key]
            else:
                merged[key] = custom[key]
        for key in defaults.keys():
            if key not in custom:
                merged[key] = defaults[key]
        self._data = merged
        self._root = AttributeDictReader(self._data, filename=config_path or defaults)

    def __getattr__(self, attr):
        return getattr(self._root, attr)


class AttributeDictReader(Mapping):

    def __init__(self, dictionary, filename='', path=tuple()):
        self._dictionary = dictionary
        self._path = path
        self._filename = filename

    def __getattr__(self, attr):
        path = self._path
        try:
            value = self._dictionary[attr]
        except KeyError:
            raise AttributeError(f'no such attribute: {self._filename}:{".".join(path + (attr,))}')
        if isinstance(value, Mapping):
            return self.__class__(value, filename=self._filename, path=path+(attr,))
        else:
            return value

    def __getitem__(self, key):
        return self._dictionary[key]

    def __iter__(self):
        return iter(self._dictionary)

    def __len__(self):
        return len(self._dictionary)
