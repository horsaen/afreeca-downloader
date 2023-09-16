import os
import requests
import platform
from urllib.parse import urljoin, urlparse
from tools.formatBytes import format_bytes

def getPlaylistInfo(link):
  split_link = link.rsplit('/')
  station_no = split_link[9]

  url = "https://stbbs.afreecatv.com/api/video/get_clip_video_info.php"
  payload = "broad_no=" + station_no + "&nCheck=CHK"
  headers = {
    "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/117.0",
    "Accept": "application/json, text/javascript, */*; q=0.01",
    "Accept-Language": "en-US,en;q=0.5",
    "Accept-Encoding": "gzip, deflate, br",
    "Referer": "https://stbbs.afreecatv.com/vodclip/index.php",
    "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
    "X-Requested-With": "XMLHttpRequest",
    "Origin": "https://stbbs.afreecatv.com",
    "DNT": "1",
    "Alt-Used": "stbbs.afreecatv.com",
    "Connection": "keep-alive",
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-origin",
    "Sec-GPC": "1",
    "TE": "trailers"
  }
  res = requests.request("POST", url, data=payload, headers=headers)
  bj_id = res.json()['bj_id']
  if platform.system() == 'Windows':
    output_file = res.json()['bj_id'] + '-' + res.json()['broad_no'] + '-' + res.json()['file_start'].replace(' ', '_').replace(':', '-') + '.ts'
  output_file = res.json()['bj_id'] + '-' + res.json()['broad_no'] + '-' + res.json()['file_start'].replace(' ', '_') + '.ts'
  return bj_id, output_file

def getMasterPlaylist(link):
  response = requests.get(link)
  playlist_content = response.text
  playlists = [
    line.strip() for line in playlist_content.splitlines() if line.endswith('.m3u8')
  ]
  return playlists[-1]

def download_m3u8(link):
  master_playlist = getMasterPlaylist(link)
  base_url = "https://vod-archive-kr-cdn-z01.afreecatv.com" + master_playlist.rsplit('/', 1)[0] + '/'
  response = requests.get("https://vod-archive-kr-cdn-z01.afreecatv.com" + master_playlist)
  playlist_content = response.text

  bj_id, output_filename = getPlaylistInfo(link)

  output_path = 'downloads/' + bj_id + '/' + output_filename

  if os.path.exists('downloads/' + bj_id) is False:
    os.makedirs('downloads/' + bj_id)

  segments = [
    line.strip() for line in playlist_content.splitlines() if line.endswith('.ts')
  ]

  file_size = 0
  current_segment = 0

  with open(output_path, 'ab') as output_file:
    for segment in segments:
      segment_url = urljoin(base_url, segment)
      response = requests.get(segment_url)
      file_size += len(response.content)
      output_file.write(response.content)
      current_segment += 1
      print("\r" + f"Downloading to {output_filename} || {format_bytes(file_size)} @ Segment {current_segment} of {len(segments)}   \x1b[?25l", end="", flush=True)