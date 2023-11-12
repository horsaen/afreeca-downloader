import requests
import os
import platform
import time
from urllib.parse import urljoin
from requests.exceptions import ReadTimeout, ConnectionError

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

from plugins.bigo.verify import checkExists, verify

def getPlaylist(id):
  url = "https://ta.bigo.tv/official_website/studio/getInternalStudioInfo"

  payload = "siteId=" + id + "&=verify%3D"
  headers = {
      "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0",
      "Accept": "application/json, text/plain, */*",
      "Accept-Language": "en-US,en;q=0.5",
      "Accept-Encoding": "gzip, deflate, br",
      "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
      "Origin": "https://www.bigo.tv",
      "DNT": "1",
      "Connection": "keep-alive",
      "Referer": "https://www.bigo.tv/",
      "Sec-Fetch-Dest": "empty",
      "Sec-Fetch-Mode": "cors",
      "Sec-Fetch-Site": "same-site",
      "Sec-GPC": "1",
      "TE": "trailers"
  }

  response = requests.request("POST", url, data=payload, headers=headers)

  return response.json()['data']['hls_src'], response.json()['data']['siteId'], response.json()['data']['nick_name']

# this sometimes drops packets, there's literally nothing i can do about it -- server issue
def downloadStream(url, siteId, nickname):
  segment_urls = set()

  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
  if platform.system() == 'Windows':
    now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

  output_filename = nickname + '-' + siteId + '-' + now + '-bigo.ts'
  output_path = 'downloads/Bigo/' + nickname + '/' + output_filename

  if os.path.exists('downloads/Bigo/' + nickname) == False:
    os.makedirs('downloads/Bigo/' + nickname)

  file_size = 0
  start_time = time.time()

  while True:
    try:
      base_url = url.rsplit('/', 1)[0] + '/'
      res = requests.get(url)
      playlist_content = res.text
      
      if '.ts' not in playlist_content:
        if verify(siteId):
          url, siteId, nickname = getPlaylist(siteId)
          downloadStream(url, siteId, nickname)

      new_segment_lines = [
        line.strip() for line in playlist_content.splitlines() if line.endswith('.ts')
      ]

      with open(output_path, 'ab') as output_file:
        for new_segment_line in new_segment_lines:
          segment_url = urljoin(base_url, new_segment_line)
          if segment_url not in segment_urls:
            # try:
            res = requests.get(segment_url, timeout=60)
            segment_urls.add(segment_url)
            # except (ReadTimeout, ConnectionError):
              # continue
            file_size += len(res.content)
            elapsed_time = time.time() - start_time
            output_file.write(res.content)
            print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)
      
    except (ReadTimeout, ConnectionError):
      continue

def main(id):
  if checkExists(id) == False:
    print('Streamer not found.')
    exit(1)
  url, siteId, nickname = verify(id)
  downloadStream(url, siteId, nickname)