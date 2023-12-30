import requests, re, time, platform, os
from urllib.parse import urljoin

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

def mainProcess(username):
  # room_id = verify(username)
  # playlist = getStreamPlaylist(room_id)
  # downloadPlaylist(playlist, username)
  getRoomId(username)

##########################################################################
############################# DOWNLOAD STREAM ############################
##########################################################################

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

def getRoomId(username):
  res = requests.get('https://www.tiktok.com/@' + username +'/live', allow_redirects=False)
  content = res.text
  # pattern = r'room_id=(\d{19})"'
  with open("res.txt", ab) as file:
    file.write(content)
  # # Search for the pattern in the HTML file
  # match = re.search(pattern, content)
  
  # if match:
  #   room_id = match.group(1)
  #   print('Room ID: ' + room_id)
  #   return room_id
  # else:
  #   print('Pattern not found.')

def getStreamPlaylist(room_id):
  res = requests.get('https://www.tiktok.com/api/live/detail/?aid=1988&roomID=' + room_id)
  content = res.json()['LiveRoomInfo']['liveUrl']
  return content

def downloadPlaylist(url, username):
  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
  if platform.system() == 'Windows':
    now = now.replace(':', '-')

  output_filename = username + '-' + now + '-tiktok.ts'
  output_path = 'downloads/TikTok/' + username + '/' + output_filename

  if os.path.exists('downloads/TikTok/' + username) == False:
    os.makedirs('downloads/TikTok/' + username)

  segment_urls = set()

  file_size = 0
  start_time = time.time()

  while True:
    res = requests.get(url)
    content = res.text
    
    if '.ts' not in content:
      downloadPlaylist(verify(username), username)

    base_url = url.rsplit('/', 1)[0] + '/'

    segments = [
      line.strip() for line in content.splitlines() if line.__contains__('.ts')
    ]

    with open(output_path, 'ab') as output_file:
      for segment in segments:
        segment_url = urljoin(base_url, segment)
        if segment_url not in segment_urls:
          segment_urls.add(segment_url)
          res = requests.get(segment_url)
          file_size += len(res.content)
          elapsed_time = time.time() - start_time
          output_file.write(res.content)
          print("\r" + f"Downloading to {output_filename} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)