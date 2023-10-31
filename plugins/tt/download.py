import requests
import time
from urllib.parse import urljoin
import platform
import os

from plugins.tt.verify import verify
from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

def downloadPlaylist(url, username):
  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
  if platform.system() == 'Windows':
    now = now.replace(':', '-')

  output_filename = username + '-' + now + '-tiktok.ts'
  output_path = 'downloads/TikTok/' + username + '/' + output_filename

  if os.path.exists('downloads/TikTok/' + username) == False:
    os.makedirs('downloads/TikTok/' + username)

  segment_urls = set()

  file_size = 0
  start_time = time.time()

  while True:
    res = requests.get(url)
    content = res.text
    
    if '.ts' not in content:
      downloadPlaylist(verify(username), username)

    base_url = url.rsplit('/', 1)[0] + '/'

    segments = [
      line.strip() for line in content.splitlines() if line.__contains__('.ts')
    ]

    with open(output_path, 'ab') as output_file:
      for segment in segments:
        segment_url = urljoin(base_url, segment)
        if segment_url not in segment_urls:
          segment_urls.add(segment_url)
          res = requests.get(segment_url)
          file_size += len(res.content)
          elapsed_time = time.time() - start_time
          output_file.write(res.content)
          print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)