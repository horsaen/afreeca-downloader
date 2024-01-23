import requests, time, os, platform
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

def mainProcess(id):
  if verify(id):
    download(id)

##########################################################################
############################# DOWNLOAD STREAM ############################
##########################################################################
  
def verify(id):
  cookie = open('cookies/flex', 'r').read().strip()
  url = "https://api.flextv.co.kr/api/channels/" + id + "/stream?option=all"

  headers = {
    'Cookie': 'flx_oauth_access=' + cookie + ';'
  }

  while True:
    res = requests.get(url, headers=headers)

    try:
      if res.json()['message'] == "방송자가 방송을 종료 했습니다.":
        print('Streamer is offline, rechecking in three minutes.')
        time.sleep(180)
    except:
      return True

def getPlaylist(id):
  cookie = open('cookies/flex', 'r').read().strip()
  url = "https://api.flextv.co.kr/api/channels/" + id + "/stream?option=all"

  headers = {
    'Cookie': 'flx_oauth_access=' + cookie + ';'
  }

  res = requests.get(url, headers=headers)

  return res.json()['sources'][0]['url'], str(res.json()['id']), res.json()['owner']['loginId']

def getStream(url):
  res = requests.get(url)

  streams = []

  for line in res.text.splitlines():  
    if line.startswith("https://"):
      streams.append(line)

  return streams[0]

def download(id):

  session = requests.Session()
  retry = Retry(connect=5, backoff_factor=0.5)
  adapter = HTTPAdapter(max_retries=retry)
  session.mount('http://', adapter)
  session.mount('https://', adapter)

  url, userId, username = getPlaylist(id)
  url = getStream(url)

  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
    
  if platform.system() == 'Windows':
    now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

  output_filename = username + '-' + userId + '-' + now + '-flex.ts'

  output_path = 'downloads/Flex/' + username + '/' + output_filename

  if os.path.exists('downloads/Flex/' + username) is False:
    os.makedirs('downloads/Flex/' + username)

  segment_urls = set()

  file_size = 0
  start_time = time.time()

  while True:
    response = session.get(url)
    playlist_content = response.text

    if '.ts' not in playlist_content:
      if verify(id):
        download(id)

    new_segment_lines = [
        line.strip() for line in playlist_content.splitlines() if line.startswith("https://")
    ]

    with open(output_path, 'ab') as output_file:
        
        for new_segment_line in new_segment_lines:
          if new_segment_line not in segment_urls:
            segment_urls.add(new_segment_line)
            response = session.get(new_segment_line)
            file_size += len(response.content)
            elapsed_time = time.time() - start_time
            output_file.write(response.content)
            print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)

    time.sleep(3)