'''
Integer based timestamps
'''

from datetime import datetime, timezone

def now():
    '''Integer based current timestamp (Unix epoch)'''
    return int(datetime.now(tz=timezone.utc).timestamp())
