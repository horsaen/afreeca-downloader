import platform
import requests
import time
import os
from urllib.parse import urljoin

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

# idk what this is but i found it in a gh repo and it looked important so here it is
clientId = "kimne78kx3ncx6brgo4mv6wki5h1ko";

def getAccessToken(username):
  url = "https://gql.twitch.tv/gql"
  payload = "{\"operationName\":\"PlaybackAccessToken_Template\",\"query\":\"query PlaybackAccessToken_Template($login: String!, $isLive: Boolean!, $vodID: ID!, $isVod: Boolean!, $playerType: String!) {  streamPlaybackAccessToken(channelName: $login, params: {platform: \\\"web\\\", playerBackend: \\\"mediaplayer\\\", playerType: $playerType}) @include(if: $isLive) {    value    signature   authorization { isForbidden forbiddenReasonCode }   __typename  }  videoPlaybackAccessToken(id: $vodID, params: {platform: \\\"web\\\", playerBackend: \\\"mediaplayer\\\", playerType: $playerType}) @include(if: $isVod) {    value    signature   __typename  }}\",\"variables\":{\"isLive\":true,\"login\":\"" + username + "\",\"isVod\":false,\"vodID\":\"\",\"playerType\":\"site\"}}"

  headers = {
    "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
    "Accept": "*/*",
    "Accept-Language": "en-US",
    "Accept-Encoding": "gzip, deflate, br",
    "Referer": "https://www.twitch.tv/",
    "Authorization": "undefined",
    "Client-ID": clientId,
    "Content-Type": "text/plain; charset=UTF-8",
    "Origin": "https://www.twitch.tv",
    "DNT": "1",
    "Connection": "keep-alive",
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-site",
    "Sec-GPC": "1"
  }

  res = requests.request("POST", url, data=payload, headers=headers)

  if res.json()['data']['streamPlaybackAccessToken'] == None:
    print("User not found")
    exit()
  else:
    return res.json()['data']['streamPlaybackAccessToken']['value'], res.json()['data']['streamPlaybackAccessToken']['signature']

def getMasterPlaylist(username, value, signature):
  params = {
    "allow_source": "true",
    "allow_audio_only": "true",
    "client_id": clientId,
    "token": value,
    "sig": signature
  }

  url = "https://usher.ttvnw.net/api/channel/hls/" + username + ".m3u8"

  while True:
    res = requests.request("GET", url, params=params)

    if res.status_code == 200:
      streams = []

      for lines in res.text.splitlines():
        if lines.startswith('https://'):
          streams.append(lines)

      return streams[0]
    
    print('Streamer is offline, rechecking in three minutes.')
    time.sleep(180)

def download(username, url):
  res = requests.request("GET", url)

  base_url = url.rsplit('/', 1)[0] + '/'
  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())

  if platform.system() == 'Windows':
    now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

  output_filename = username + '-' + now + '-twitch.ts'

  output_path = 'downloads/Twitch/' + username + '/' + output_filename

  if os.path.exists('downloads/Twitch/' + username) is False:
    os.makedirs('downloads/Twitch/' + username)

  segment_urls = set()

  file_size = 0
  start_time = time.time()

  while True:
    res = requests.get(url)
    playlist_content = res.text

    if '.ts' or '.TS' not in playlist_content:
      value, signature = getAccessToken(username)
      newUrl = getMasterPlaylist(username, value, signature)
      download(username, newUrl)

    new_segment_lines = [
      line.strip() for line in playlist_content.splitlines() if line.endswith('.ts') or line.endswith('.TS')
    ]

    with open(output_path, 'ab') as output_file:
      for new_segment_line in new_segment_lines:
        segment_url = urljoin(base_url, new_segment_line)
        if segment_url not in segment_urls:
          segment_urls.add(segment_url)
          res = requests.get(segment_url)
          file_size += len(res.content)
          elapsed_time = time.time() - start_time
          output_file.write(res.content)
          print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)

def main(username):
  value, signature = getAccessToken(username)
  url = getMasterPlaylist(username, value, signature)

  download(username, url)