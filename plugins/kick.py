import requests, time, platform, os
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

def getUrl(username):

  

  url = "https://kick.com/api/v1/channels/" + username

  cookies = []
  with open('cookies/kick', 'r') as f:
    for line in f:
      cookies.append(line.strip())

  payload = ""
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

  response = requests.request("GET", url, data=payload, headers=headers)

  return response.json()['playback_url']

def getStreams(username):
  url = getUrl(username)

  streams = requests.get(url)

  stream = []

  for lines in streams.text.splitlines():
    if lines.startswith("https://"):
      stream.append(lines)

  return stream[0]

def downloadStream(username):

  # attempt to alleviate dropped connections [NEEDS TESTING]
  session = requests.Session()
  retry = Retry(connect=5, backoff_factor=0.5)
  adapter = HTTPAdapter(max_retries=retry)
  session.mount('http://', adapter)
  session.mount('https://', adapter)

  url = getStreams(username)

  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
    
  if platform.system() == 'Windows':
    now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

  output_filename = username + '-' + now + '-kick.ts'

  output_path = 'downloads/Kick/' + username + '/' + output_filename

  if os.path.exists('downloads/Kick/' + username) is False:
    os.makedirs('downloads/Kick/' + username)

  segment_urls = set()

  file_size = 0
  start_time = time.time()

  while True:
    response = requests.get(url)
    playlist_content = response.text

    new_segment_lines = [
      line.strip() for line in playlist_content.splitlines() if line.startswith("https://")
    ]

    with open(output_path, 'ab') as output_file:
        for new_segment_line in new_segment_lines:
          segment_url = new_segment_line
          if segment_url not in segment_urls:
            segment_urls.add(segment_url)
            response = session.get(segment_url)
            file_size += len(response.content)
            elapsed_time = time.time() - start_time
            output_file.write(response.content)
            print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)

    time.sleep(3)