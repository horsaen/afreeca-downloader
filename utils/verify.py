import requests
import time
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

def verify(username):
    session = requests.Session()
    retry = Retry(connect=3, backoff_factor=0.5)
    adapter = HTTPAdapter(max_retries=retry)
    session.mount('http://', adapter)
    session.mount('https://', adapter)

    # standalone headers as api doesn't care
    headers = {'user-agent': 'Mozilla/5.0 (X11; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0'}

    while True:
        res = session.get('https://bjapi.afreecatv.com/api/' + username + '/station', headers=headers)
        if res.json().get('code') is None:
            if res.json()['broad'] is not None:
                print('Streamer is online.')
                return 1
        else:
            print('Streamer not found.')
            exit(1)

        print('Streamer is offline, rechecking in three minutes.')
        time.sleep(180)