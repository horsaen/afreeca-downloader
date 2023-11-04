import requests
import time
import os
import platform
from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration
from plugins.pandatv.verify import verify
from requests.exceptions import ReadTimeout, ConnectionError

def download(url, username):

  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())

  if platform.system() == 'Windows':
    now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

  output_filename = username + '-' + now + '-pandatv.ts'
  output_path = 'downloads/PandaTV/' + username + '/' + output_filename

  if os.path.exists('downloads/PandaTV/' + username) == False:
    os.makedirs('downloads/PandaTV/' + username)

  segment_urls = set()

  file_size = 0
  start_time = time.time()

  while True:
    try:
      response = requests.get(url)
      playlist_content = response.text

      if '.ts' not in playlist_content:
        download(verify(username), username)

      new_segment_lines = [
        line.strip() for line in playlist_content.splitlines() if line.startswith('https://')
      ]
      
      with open(output_path, 'ab') as output_file:
        for new_segment_line in new_segment_lines:
          segment_url = new_segment_line
          if segment_url not in segment_urls:
            segment_urls.add(segment_url)
            try:
              res = requests.get(segment_url, timeout=10)
            except (ReadTimeout, ConnectionError):
              continue
            file_size += len(res.content)
            elapsed_time = time.time() - start_time
            output_file.write(res.content)
            print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)

      time.sleep(1)

    except (ReadTimeout, ConnectionError):
      continue