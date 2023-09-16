import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
import time
from tools.formatBytes import format_bytes
from utils.getPlaylist import getStationNo
import os

def downloadVod(url, output_filename, username):
   
    # attempt to remedy dropped connections, works i think ???/
    session = requests.Session()
    retry = Retry(connect=3, backoff_factor=0.5)
    adapter = HTTPAdapter(max_retries=retry)
    session.mount('http://', adapter)
    session.mount('https://', adapter)

    output_path = 'downloads/' + username + '/' + output_filename

    if os.path.exists('downloads/' + username) is False:
        os.makedirs('downloads/' + username)

    while True:
        headers = {}
        if os.path.exists(output_path):
            headers['Range'] = f'bytes={os.path.getsize(output_path)}-'

        res = session.get(url, stream=True)

        with open(output_path, 'ab') as output_file:
            for chunk in res.iter_content(chunk_size=1024):
                output_file.write(chunk)
                # does not include time as it is not a live stream, don't want to run mediainfo on every check
                if len(chunk) != 1024 and getStationNo(username, '') is False:
                    print('\nDownload complete.')
                    break
                    
                print("\r" + f"Downloading to {output_filename} || {format_bytes(os.path.getsize(output_path))}    \x1b[?25l", end="", flush=True)

        time.sleep(1)
        exit()