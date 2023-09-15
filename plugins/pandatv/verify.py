import requests
import time
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

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