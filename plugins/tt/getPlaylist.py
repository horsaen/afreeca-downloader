import requests
import re

def getRoomId(username):
  res = requests.get('https://www.tiktok.com/@' + username +'/live', allow_redirects=False)
  content = res.text
  pattern = r'room_id=(\d{19})"'

  # Search for the pattern in the HTML file
  match = re.search(pattern, content)
  
  if match:
    room_id = match.group(1)
    print('Room ID: ' + room_id)
    return room_id
  else:
    print('Pattern not found.')

def getStreamPlaylist(room_id):
  res = requests.get('https://www.tiktok.com/api/live/detail/?aid=1988&roomID=' + room_id)
  content = res.json()['LiveRoomInfo']['liveUrl']
  return content