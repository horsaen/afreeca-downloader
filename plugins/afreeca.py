import platform, time, requests, os, re
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from urllib.parse import urljoin

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

def mainProcess(mode, username, password):
  if mode == "fromStart":
    station_no = getStationNo(username, password)
    if station_no is False:
      station_no = input("Unable to get get stream id automatically, please input one manually:\n")
    vodUrl, filename = getVodPlaylist(username, station_no)
    downloadVod(vodUrl, filename, username)
  elif mode == "playlist":
    downloadPlaylists()
  elif mode == "live":
    if verify(username):
      download(getVideoPlaylist(username, password), username)

##########################################################################
############################# DOWNLOAD STREAM ############################
##########################################################################

def verify(username):
  session = requests.Session()
  retry = Retry(connect=3, backoff_factor=0.5)
  adapter = HTTPAdapter(max_retries=retry)
  session.mount('http://', adapter)
  session.mount('https://', adapter)

  # standalone headers as api doesn't care
  headers = {'user-agent': 'Mozilla/5.0 (X11; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0'}

  while True:
    res = session.get('https://bjapi.afreecatv.com/api/' + username + '/station', headers=headers)
    if res.json().get('code') is None:
      if res.json()['broad'] is not None:
        print('Streamer is online.')
        return 1
    else:
      print('Streamer not found.')
      exit(1)

    print('Streamer is offline, rechecking in three minutes.')
    time.sleep(180)

def getPlaylistBase(station_no):
  url = 'https://livestream-manager.afreecatv.com/broad_stream_assign.html?return_type=gcp_cdn&cors_origin_url=play.afreecatv.com&broad_key=' + station_no + '-common-master-hls'
  # this one doesn't even need headers, yay :)
  res = requests.get(url)
  return res.json()['view_url']

def getStationNo(username, pwd):
  url = 'https://bjapi.afreecatv.com/api/' + username + '/station'
  # this one just needs simple headers
  headers = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0",
    "Accept": "application/json, text/plain, */*"
  }
  res = requests.get(url, headers=headers)
  if res.json()['broad'] is None:
    return False
  else:
    return str(res.json()['broad']['broad_no'])
  
def getMasterPlaylist(username, station_no, pwd):
  url = "https://live.afreecatv.com/afreeca/player_live_api.php"
  cookie = open('cookies/afreeca', 'r').read().strip()
  payload = "bid=" + username + "&bno=" + station_no + "&type=aid&pwd=" + pwd + "&player_type=html5&stream_type=common&quality=master&mode=landing&from_api=0"
  # this one actively needs a cookie for 19+ streams
  headers = {
    "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/118.0",
    "Accept": "*/*",
    "Accept-Language": "en-US,en;q=0.5",
    "Accept-Encoding": "gzip, deflate, br",
    "Content-Type": "application/x-www-form-urlencoded",
    "Origin": "https://play.afreecatv.com",
    "DNT": "1",
    "Alt-Used": "live.afreecatv.com",
    "Connection": "keep-alive",
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-site",
    "Sec-GPC": "1",
    "Cookie": "PdboxTicket=" + cookie,
    "TE": "trailers",
  }
  res = requests.request("POST", url, data=payload, headers=headers)
  return res.json()['CHANNEL']['AID']

def getStreamList(base, master):
  url = base + "?aid=" + master
  res = requests.get(url)
  streams = []
  for lines in res.text.splitlines():
    if lines.startswith('auth_playlist'):
      streams.append(lines)
  # 0, original, 1920x1080
  # 1, hd, 960x540
  # 2, sd, 640x360
  return streams[0]

def getVideoPlaylist(username, pwd):
  try:
    base = getPlaylistBase(getStationNo(username, ''))
    master = getMasterPlaylist(username, getStationNo(username, ''), '')
    playlist = getStreamList(base, master)
  except KeyError:
    base = getPlaylistBase(getStationNo(username, ''))
    master = getMasterPlaylist(username, getStationNo(username, ''), '')
    playlist = getStreamList(base, master)

  return str(base + '/' + playlist)

def download(url, username):
  # attempt to alleviate dropped connections [NEEDS TESTING]
  session = requests.Session()
  retry = Retry(connect=5, backoff_factor=0.5)
  adapter = HTTPAdapter(max_retries=retry)
  session.mount('http://', adapter)
  session.mount('https://', adapter)
  cookie = open('cookies/afreeca', 'r').read().strip()
  headers = {
    "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/118.0",
    "Accept": "*/*",
    "Accept-Language": "en-US,en;q=0.5",
    "Accept-Encoding": "gzip, deflate, br",
    "Content-Type": "application/x-www-form-urlencoded",
    "Origin": "https://play.afreecatv.com",
    "DNT": "1",
    "Alt-Used": "live.afreecatv.com",
    "Connection": "keep-alive",
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-site",
    "Sec-GPC": "1",
    "Cookie": cookie,
    "TE": "trailers",
  }

  base_url = url.rsplit('/', 1)[0] + '/'
  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
  
  if platform.system() == 'Windows':
    now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

  output_filename = username + '-' + getStationNo(username,'') + '-' + now + '-afreeca.ts'

  output_path = 'downloads/Afreeca/' + username + '/' + output_filename

  if os.path.exists('downloads/Afreeca/' + username) is False:
      os.makedirs('downloads/Afreeca/' + username)

  segment_urls = set()

  file_size = 0
  start_time = time.time()

  while True:
    response = session.get(url, headers=headers)
    playlist_content = response.text

    if '<p>Sorry, the page you are looking for is currently unavailable.<br/>' in playlist_content:
        print("\nStream paused")
        time.sleep(60)
        continue

    # works for the most part, but needs testing
    if '<title>502 Server Error</title>' in playlist_content:
        print(f"\nStream finished. Exiting the program.")
        if verify(username):
            download(getVideoPlaylist(username, ''), username)

    new_segment_lines = [
        line.strip() for line in playlist_content.splitlines() if line.endswith('.TS') or line.endswith('.ts')
    ]

    with open(output_path, 'ab') as output_file:
        
        for new_segment_line in new_segment_lines:
            segment_url = urljoin(base_url, new_segment_line)
            if segment_url not in segment_urls:
                segment_urls.add(segment_url)
                response = session.get(segment_url)
                file_size += len(response.content)
                elapsed_time = time.time() - start_time
                output_file.write(response.content)
                print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)

    time.sleep(3)

##########################################################################
############################ CONCURRENT STREAMS ##########################
##########################################################################
    
def getUserData(username):
  url = 'https://bjapi.afreecatv.com/api/' + username + '/station'
  # this one just needs simple headers
  headers = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0",
    "Accept": "application/json, text/plain, */*"
  }
  res = requests.get(url, headers=headers)
  return res.json()['station']['user_id'], res.json()['station']['user_nick'], str(res.json()['broad']['broad_no'])

def concurrentVerify(username):
  session = requests.Session()
  retry = Retry(connect=3, backoff_factor=0.5)
  adapter = HTTPAdapter(max_retries=retry)
  session.mount('http://', adapter)
  session.mount('https://', adapter)

  # standalone headers as api doesn't care
  headers = {'user-agent': 'Mozilla/5.0 (X11; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0'}

  while True:
    res = session.get('https://bjapi.afreecatv.com/api/' + username + '/station', headers=headers)
    if res.json().get('code') is None:
      if res.json()['broad'] is not None:
        return True

    time.sleep(180)

##########################################################################
############################## DOWNLOAD VOD ##############################
##########################################################################

def getVodPlaylist(username, station_no):
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
  try:
    vodUrl = 'https://vod-archive-fmp4-kr-cdn-z01.afreecatv.com' + res.json()['media_path']
    if platform.system() == 'Windows':
      filename = res.json()['bj_id'] + '-' + res.json()['broad_no'] + '-' + res.json()['file_start'].replace(' ', '_').replace(':', '-') + '.mp4'
    filename = res.json()['bj_id'] + '-' + res.json()['broad_no'] + '-' + res.json()['file_start'].replace(' ', '_') + '.mp4'
  except KeyError:
    exit('Error getting VOD playlist.\nCould be a wrong station number or the stream has expired.')
  return vodUrl, filename

def downloadVod(url, output_filename, username):
   
  # attempt to remedy dropped connections, works i think ???/
  session = requests.Session()
  retry = Retry(connect=3, backoff_factor=0.5)
  adapter = HTTPAdapter(max_retries=retry)
  session.mount('http://', adapter)
  session.mount('https://', adapter)

  output_path = 'downloads/Afreeca/' + username + '/' + output_filename

  if os.path.exists('downloads/Afreeca/' + username) is False:
      os.makedirs('downloads/Afreeca/' + username)

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
                  exit()
                  
              print("\r" + f"Downloading to {output_filename} || {format_bytes(os.path.getsize(output_path))}    \x1b[?25l", end="", flush=True)

      time.sleep(1)

##########################################################################
############################## DOWNLOAD M3U8 #############################
##########################################################################

def downloadPlaylists():
  inputStreams = []

  print("Paste playlist URLS:")

  while True:
    user_input = input("")

    if user_input == "":
      print("Downloading vod playlists...")
      break
    else:
      inputStreams.append(user_input)

  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
    
  if platform.system() == 'Windows':
    now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

  for i, stream in enumerate(inputStreams):
    streams = []

    base_url = re.match(r"(https?://[^/]+)", stream).group(1)

    res = requests.get(stream)

    for lines in res.text.splitlines():
      if lines.endswith("playlist.m3u8"):
        streams.append(lines)
      
    url = base_url + streams[-1]
    res = requests.get(url)

    playlistBase = url.rsplit('/', 1)[0] + '/'

    vodSegments = [
      line.strip() for line in res.text.splitlines() if line.endswith('.ts')
    ]

    segmentSize = 0

    output_filename = str(i+1) + ".ts"

    if os.path.exists("downloads/Afreeca/m3u8/" + now) is False:
      os.makedirs("downloads/Afreeca/m3u8/" + now)

    output_path = "downloads/Afreeca/m3u8/" + now + "/" + output_filename

    with open(output_path, 'ab') as output_file:
      for segment in vodSegments:
        segmentSize += 1
        response = requests.get(playlistBase + segment)
        output_file.write(response.content)
        print("\r" + f"Downloading to {now}/{output_filename} || Segment {str(segmentSize-1)}/{str(len(vodSegments)-1)} || Vod {i+1} of {len(inputStreams)}      \x1b[?25l", end="", flush=True)