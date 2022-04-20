'''
Integer based timestamps
'''

from datetime import datetime, timezone
from dateutil.parser import parse

def now():
    '''Integer based current timestamp (Unix epoch)'''
    return int(datetime.now(tz=timezone.utc).timestamp())


def from_string(string):
    '''Read timestamp from string'''
    dtime = parse(string)
    return int(dtime.astimezone(timezone.utc).timestamp())
