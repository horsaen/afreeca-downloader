import requests
import time

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
  "Cookie": "sessKey=" + cookie,
  "Sec-Fetch-Dest": "empty",
  "Sec-Fetch-Mode": "cors",
  "Sec-Fetch-Site": "same-site",
  "Sec-GPC": "1",
  "TE": "trailers"
}

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
    except:
      # do this all here instead of needing another file
      response = requests.get(res.json()['PlayList']['hls'][0]['url'])

      steams = []

      for lines in response.text.splitlines():
        if lines.startswith('https://'):
          steams.append(lines)

      return steams[0]

    time.sleep(180)
