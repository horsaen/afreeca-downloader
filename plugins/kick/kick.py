import requests
import platform
import time
import os
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

def getPlaylist(username):
  url = "https://kick.com/api/v2/channels/" + username + "/livestream"

  cookies = []
  with open('cookies/kick', 'r') as f:
    for line in f:
      cookies.append(line.strip())

  headers = {
    "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Accept": "application/json, text/plain, */*",
    "Accept-Language": "en-US",
    "Accept-Encoding": "gzip, deflate, br",
    "X-Socket-ID": "95721.13970",
    "DNT": "1",
    "Connection": "keep-alive",
    'Cookie': '__cf_bm=' + cookies[0] + '; cf_clearance=' + cookies[1] + '; kick_session=' + cookies[2] + ';',
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-origin",
    "Sec-GPC": "1",
    "TE": "trailers"
  }

  response = requests.request("GET", url, headers=headers)

  return response.json()['data']['playback_url']

def getVideo(username):
  url = getPlaylist(username)

  res = requests.get(url)

  streams = []

  for lines in res.text.splitlines():
    if lines.startswith('https://'):
      streams.append(lines)

  return streams[0]

def download(username):
  session = requests.Session()
  retry = Retry(connect=3, backoff_factor=0.5)
  adapter = HTTPAdapter(max_retries=retry)
  session.mount('http://', adapter)
  session.mount('https://', adapter)

  url = getVideo(username)

  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
  if platform.system() == 'Windows':
    now = now.replace(':', '-')
  output_filename = username + '-' + now + '-kick.ts'
  output_path = 'downloads/Kick/' + username + '/' + output_filename

  if os.path.exists('downloads/Kick/' + username) == False:
    os.makedirs('downloads/Kick/' + username)

  segment_urls = set()

  file_size = 0
  start_time = time.time()

  while True:

    response = session.get(url)
    playlist_content = response.text

    new_segment_lines = [
      line.strip() for line in playlist_content.splitlines() if line.startswith('https://')
    ]
    
    with open(output_path, 'ab') as output_file:
      for new_segment_line in new_segment_lines:
        segment_url = new_segment_line
        if segment_url not in segment_urls:
          segment_urls.add(segment_url)
          res = requests.get(segment_url)
          file_size += len(res.content)
          elapsed_time = time.time() - start_time
          output_file.write(res.content)
          print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)