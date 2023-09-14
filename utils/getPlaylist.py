import requests

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
  cookie = open('cookies', 'r').read().strip()
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
    "Cookie": cookie,
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
    filename = res.json()['bj_id'] + '-' + res.json()['broad_no'] + '-' + res.json()['file_start'].replace(' ', '_') + '.mp4'
  except KeyError:
    exit('Error getting VOD playlist.\nCould be a wrong station number or the stream has expired.')
  return vodUrl, filename

def getVideoPlaylist(username, pwd):
  base = getPlaylistBase(getStationNo(username, ''))
  master = getMasterPlaylist(username, getStationNo(username, ''), '')
  playlist = getStreamList(base, master)

  return str(base + '/' + playlist)