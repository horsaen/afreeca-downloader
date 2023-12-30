import requests, time, os, platform
from requests.exceptions import ReadTimeout, ConnectionError

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

cookie = open('cookies/panda', 'r').read().strip()
headers = {
  "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/117.0",
  "Accept": "application/json, text/plain, */*",
  "Accept-Language": "en-US,en;q=0.5",
  "Accept-Encoding": "gzip, deflate, br",
  "Content-Type": "application/x-www-form-urlencoded",
  "Origin": "https://www.pandalive.co.kr",
  "DNT": "1",
  "Connection": "keep-alive",
  "Referer": "https://www.pandalive.co.kr/",
  "Cookie": "sessKey=" + cookie,
  "Sec-Fetch-Dest": "empty",
  "Sec-Fetch-Mode": "cors",
  "Sec-Fetch-Site": "same-site",
  "Sec-GPC": "1",
  "TE": "trailers"
}

def mainProcess(username):
  if verify(username):
    if os.path.exists('downloads/PandaTV') is False:
      os.makedirs('downloads/PandaTV')

  download(verify(username), username)

##########################################################################
############################# DOWNLOAD STREAM ############################
##########################################################################

def checkUser(username):
  url = "https://api.pandalive.co.kr/v1/member/bj"

  payload = "userId=" + username

  response = requests.request("POST", url, data=payload, headers=headers)
  
  if response.json()['result']:
    return(response.json()['result'])
  else:
    print('Streamer not found')
    exit(1)

# get master plist
def verify(username):
  url = "https://api.pandalive.co.kr/v1/live/play"

  payload = 'action=watch&userId=' + username + '&password=&width=48&height=48&imageResize=crop&fanIconWidth=44&fanIconHeight=44&fanIconImageResize=crop'

  while True:
    res = requests.request("POST", url, data=payload, headers=headers)
    # check if error
    try:
      if res.json()['errorData'] is not None:
        if res.json()['errorData']['code'] == 'needAdult':
          print('Stream is 19+ and unable to retrieve stream URL, input a valid sessKey in panda-cookies.')
          exit(1)
        elif res.json()['errorData']['code'] == 'castEnd':
          print('Stream is offline, retrying in 3 minutes.')
    # handle if no error
    except KeyError:
      # do this all here instead of needing another file
      try:
        response = requests.get(res.json()['PlayList']['hls'][0]['url'])

        steams = []

        for lines in response.text.splitlines():
          if lines.startswith('https://'):
            steams.append(lines)

        return steams[0]
      except TypeError:
        print('Unhandled error, trying again in 3 minutes.')

    time.sleep(180)



def getPlaylist(username):
  url = "https://api.pandalive.co.kr/v1/live/play"

  payload = 'action=watch&userId=' + username + '&password=&width=48&height=48&imageResize=crop&fanIconWidth=44&fanIconHeight=44&fanIconImageResize=crop'

  res = requests.request("POST", url, data=payload, headers=headers)

  response = requests.get(res.json()['PlayList']['hls'][0]['url'])

  steams = []

  for lines in response.text.splitlines():
    if lines.startswith('https://'):
      steams.append(lines)

  return steams[0]

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


##########################################################################
############################ CONCURRENT STREAMS ##########################
##########################################################################

def concurrentVerify(username):
  url = "https://api.pandalive.co.kr/v1/live/play"

  payload = 'action=watch&userId=' + username + '&password=&width=48&height=48&imageResize=crop&fanIconWidth=44&fanIconHeight=44&fanIconImageResize=crop'

  while True:
    res = requests.request("POST", url, data=payload, headers=headers)
    try:
      if res.json()['errorData'] is not None:
        if res.json()['errorData']['code'] == 'needAdult':
          return 'Err19'
        elif res.json()['errorData']['code'] == 'castEnd':
          continue
    except KeyError:
      return True
  
    time.sleep(180)

def concurrentVerifyOld(username):
  url = "https://api.pandalive.co.kr/v1/live/play"

  payload = 'action=watch&userId=' + username + '&password=&width=48&height=48&imageResize=crop&fanIconWidth=44&fanIconHeight=44&fanIconImageResize=crop'

  while True:
    res = requests.request("POST", url, data=payload, headers=headers)

    if res.json()['errorData']['code'] == 'castEnd':

      if res.json()['errorData'] is not None:
        if res.json()['errorData']['code'] == 'needAdult':
          return 'Err19'
        elif res.json()['errorData']['code'] == 'castEnd':
          continue
      else:
        return True

    time.sleep(180)