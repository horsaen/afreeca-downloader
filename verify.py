import requests
import time
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

def verify(username, url):

    # attempt to retry the request if it's refused as to not be thrown out of the program
    session = requests.Session()
    retry = Retry(connect=3, backoff_factor=0.5)
    adapter = HTTPAdapter(max_retries=retry)
    session.mount('http://', adapter)
    session.mount('https://', adapter)

    # supa dupa sneaky headers to ensure that you get a valid JSON response
    headers = {'user-agent': 'Mozilla/5.0 (X11; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0'}

    while True:

        response = session.get(url, headers=headers)

        if response.json()['broad'] is not None:
            print("Streamer is online.")
            streamStruct: str = "https://play.afreecatv.com/" + username + "/" + str(response.json()['broad']['broad_no'])
            return streamStruct
        
        print("Streamer is offline, rechecking in three minutes.")
        # set to 3 mins to avoid rate limiting >:(
        time.sleep(180)