import requests, os, platform, time
from urllib.parse import urljoin
from requests.exceptions import ReadTimeout, ConnectionError

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

def mainProcess(id):
  if checkExists(id) == False:
    print('Streamer not found.')
    exit(1)
  if verify(id) is True:
    url = getPlaylist(id)
    siteId, nickname = getStreamData(id)
  downloadStream(url, siteId, nickname)

##########################################################################
############################# DOWNLOAD STREAM ############################
##########################################################################

def checkExists(siteId):
  url = "https://ta.bigo.tv/official_website/studio/getInternalStudioInfo"

  payload = "siteId=" + siteId + "&=verify%3D"
  
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

  if response.json()['data']['sid'] is None:
    return False
  else:
    return True

def verify(siteId):
  url = "https://ta.bigo.tv/official_website/studio/getInternalStudioInfo"

  payload = "siteId=" + siteId + "&=verify%3D"

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

  while True:
    response = requests.request("POST", url, data=payload, headers=headers)
    if response.json()['data']['hls_src'] != "":
      return True
      # return response.json()['data']['hls_src'], response.json()['data']['siteId'], response.json()['data']['nick_name']

    print('Streamer is offline, rechecking in three minutes.')
    
    time.sleep(180)

def getStreamData(id):
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

  return response.json()['data']['siteId'], response.json()['data']['nick_name']

def getPlaylist(id):
  url = "https://ta.bigo.tv/official_website/studio/getInternalStudioInfo"

  payload = "siteId=loveyouarge&=verify%3D"
  headers = {
    "cookie": "www_random_gray=52",
    "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0",
    "Accept": "application/json, text/plain, */*",
    "Accept-Language": "en-US,en;q=0.5",
    "Accept-Encoding": "gzip, deflate, br",
    "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
    "Origin": "https://www.bigo.tv",
    "DNT": "1",
    "Connection": "keep-alive",
    "Referer": "https://www.bigo.tv/",
    "Cookie": "www_random_gray=96",
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-site",
    "Sec-GPC": "1",
    "TE": "trailers"
  }

  response = requests.request("POST", url, data=payload, headers=headers)

  return response.json()['data']['hls_src']


# no longer working
# def getPlaylist(id):
#   url = "https://www.bigo.tv/OInterface/getVideoParam?bigoId=" + id
#   headers = {
#     "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0"
#   }

#   res = requests.request("GET", url, headers=headers)

#   return res.json()['data']['videoSrc']

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
          url = getPlaylist(siteId)
          siteId, nickname = getStreamData(siteId)
          downloadStream(url, siteId, nickname)

      new_segment_lines = [
        line.strip() for line in playlist_content.splitlines() if line.endswith('.ts')
      ]

      with open(output_path, 'ab') as output_file:
        for new_segment_line in new_segment_lines:
          segment_url = urljoin(base_url, new_segment_line)
          if segment_url not in segment_urls:
            try:
              res = requests.get(segment_url, timeout=60)
              segment_urls.add(segment_url)
            except (ReadTimeout, ConnectionError):
              continue
            file_size += len(res.content)
            elapsed_time = time.time() - start_time
            output_file.write(res.content)
            print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)
      
    except (ReadTimeout, ConnectionError):
      continue

##########################################################################
############################ CONCURRENT STREAMS ##########################
##########################################################################

def concurrentVerify(siteId):
  url = "https://ta.bigo.tv/official_website/studio/getInternalStudioInfo"

  payload = "siteId=" + siteId + "&=verify%3D"

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

  while True:
    response = requests.request("POST", url, data=payload, headers=headers)
    if response.json()['data']['hls_src'] != "":
      return True
    
    time.sleep(180)