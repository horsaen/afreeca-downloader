import requests
import re
import time

def verify(username):
  while True:
    res = requests.get('https://www.tiktok.com/@' + username +'/live', allow_redirects=False)
    content = res.text
    pattern = r'room_id=(\d{19})"'

    match = re.search(pattern, content)
    if match:
      room_id = match.group(1)
      print('Room ID: ' + room_id)
      return room_id
    
    print('Streamer is offline, rechecking in three minutes.')
    time.sleep(180)