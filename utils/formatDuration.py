from datetime import timedelta

def format_duration(duration):
    td = timedelta(seconds=int(duration))
    return str(td)