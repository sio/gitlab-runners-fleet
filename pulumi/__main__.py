"""A Python Pulumi program"""

import pulumi

from data import InstanceParams
from instance import create

create(InstanceParams('test-instance'))
