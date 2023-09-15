import requests

def getMasterPlaylist(username, pwd):
  url = "https://api.pandalive.co.kr/v1/live/play"

  payload = 'action=watch&userId=' + username + '&password=&width=48&height=48&imageResize=crop&fanIconWidth=44&fanIconHeight=44&fanIconImageResize=crop'
  
  cookie = open('plugins/pandatv/panda-cookies', 'r').read().strip()
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
    "Cookie": cookie,
    "Sec-Fetch-Dest": "empty",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Site": "same-site",
    "Sec-GPC": "1",
    "TE": "trailers"
  }

  response = requests.request("POST", url, data=payload, headers=headers)

  return response.json()['PlayList']['hls'][0]['url']

def getStreams(username, pwd):
  url = getMasterPlaylist(username, pwd)
  res = requests.get(url)

  streams = []

  for lines in res.text.splitlines():
    if lines.startswith('https://'):
      streams.append(lines)

  return streams