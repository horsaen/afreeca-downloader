import requests
import time
import os
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from plugins.pandatv.download import download

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

def userExists(username):
  url = "https://api.pandalive.co.kr/v1/member/bj"

  payload = "userId=" + username

  response = requests.request("POST", url, data=payload, headers=headers)

  return(response.json()['result'])

def verify(username):
  session = requests.Session()
  retry = Retry(connect=3, backoff_factor=0.5)
  adapter = HTTPAdapter(max_retries=retry)
  session.mount('http://', adapter)
  session.mount('https://', adapter)

  while True:
    url = "https://api.pandalive.co.kr/v1/live/play"

    payload = "action=watch&userId=" + username + "&password="

    res = session.request("POST", url, data=payload, headers=headers)

    if res.json()['result'] is True:
      print('Streamer is online.')
      return 1
    if res.json()['message'] == '본인인증이 필요합니다..':
      print('Stream requires identity verification\nPlease provide valid login cookie in headers')
      exit()
    if res.json()['result'] is False:
      print('Streamer is offline, rechecking in three minutes.')
      time.sleep(180)

def main(username):
  if os.path.exists('downloads/PandaTV') is False:
    os.makedirs('downloads/PandaTV')

  exists = userExists(username)
  
  if exists is False:
    print('User does not exist.')
    exit()
  else:
    download(username)