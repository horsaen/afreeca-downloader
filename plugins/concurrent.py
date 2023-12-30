# likely a cleaner way to work this, works for now
# nice to have the threads separate in case something happens to api (?)
from enum import verify
import platform, time, requests, os, re
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from urllib.parse import urljoin
from threading import Thread
from tabulate import tabulate

from tools.formatBytes import format_bytes
from tools.formatDuration import format_duration

from plugins.afreeca import getVideoPlaylist, getUserData, concurrentVerify as afreecaVerify
from plugins.bigo import concurrentVerify as bigoVerify, getPlaylist as bigoPlaylist, getStreamData as bigoData
from plugins.panda import concurrentVerify as pandaVerify, getPlaylist as pandaPlaylist

usernameList = []

def main():
  users = []
  threads = []

  with open('users', 'r') as file:
    for line in file:
      uname, platform = line.strip().split(',')
      users.append([uname, platform])

    for user in users:
      name, platform = user
      instanceId = users.index(user)

      if platform == " afreeca":
        thread = Thread(target=afreeca, args=(instanceId, name))
        thread.start()
        threads.append(thread)
      if platform == " panda":
        thread = Thread(target=panda, args=(instanceId, name))
        thread.start()
        threads.append(thread)
      if platform == " bigo":
        thread = Thread(target=bigo, args=(instanceId, name))
        thread.start()
        threads.append(thread)
  
  while True:
    try:
      head = ["Site", "User", "Nick", "Size", "Duration", "Path"]
      os.system('cls' if os.name == 'nt' else 'clear')
      print(tabulate(usernameList, headers=head, tablefmt='grid') + '\x1b[?25l')
      time.sleep(2)
    except KeyboardInterrupt:
      exit()

##########################################################################
################################# DOWNLOADS ##############################
##########################################################################

##########################################################################
################################# AFREECA ################################
##########################################################################

# did my best to increase redudancy on afreeca, would just randomly die,,, not cool
def afreeca(instanceId, user):
  usernameList.insert(instanceId, ["Afreeca", user, '', '', '', ''])

  if afreecaVerify(user) is True:
  
    m3u8Url = getVideoPlaylist(user, '')
    username, nickname, station_no = getUserData(user)

    segment_urls = set()
    file_size = 0
    start_time = time.time()

    if os.path.exists('downloads/Afreeca/' + user) == False:
      os.makedirs('downloads/Afreeca/' + user)

    now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
    if platform.system() == 'Windows':
      now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

    output_filename = user + '-' + station_no + '-' + now + '-afreeca.ts'
    output_path = 'downloads/Afreeca/' + user + '/' + output_filename

    while True:
      base_url = m3u8Url.rsplit('/', 1)[0] + '/'

      res = requests.get(m3u8Url)
      playlist_content = res.text

      new_segment_lines = [
        line.strip() for line in playlist_content.splitlines() if line.endswith('.TS') or line.endswith('.ts')
      ]

      if '.ts' not in playlist_content.lower():
        try:
          usernameList[instanceId] = ["Afreeca", user, nickname, 'Offline', 'Offline', 'Offline']

          if afreecaVerify(user) is True:

            m3u8Url = getVideoPlaylist(user, '')

            segment_urls = set()
            file_size = 0
            start_time = time.time()

            if os.path.exists('downloads/Afreeca/' + user) == False:
              os.makedirs('downloads/Afreeca/' + user)

            now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
            if platform.system() == 'Windows':
              now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

            output_filename = user + '-' + station_no + '-' + now + '-afreeca.ts'
            output_path = 'downloads/Afreeca/' + user + '/' + output_filename

            continue
        except:
          usernameList[instanceId] = ["Afreeca", user, nickname, 'ERR', 'RETRYING', 'CHECK COOKIES!']

          if afreecaVerify(user) is True:

            m3u8Url = getVideoPlaylist(user, '')

            segment_urls = set()
            file_size = 0
            start_time = time.time()

            if os.path.exists('downloads/Afreeca/' + user) == False:
              os.makedirs('downloads/Afreeca/' + user)

            now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
            if platform.system() == 'Windows':
              now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

            output_filename = user + '-' + station_no + '-' + now + '-afreeca.ts'
            output_path = 'downloads/Afreeca/' + user + '/' + output_filename

            continue


      with open(output_path, 'ab') as output_file:
        for new_segment_line in new_segment_lines:
          segment_url = urljoin(base_url, new_segment_line)
          if segment_url not in segment_urls:
            segment_urls.add(segment_url)
            try:
              res = requests.get(segment_url, timeout=15)
            except (requests.ReadTimeout, ConnectionError):
              continue
            file_size += len(res.content)
            elapsed_time = time.time() - start_time
            output_file.write(res.content)
            usernameList[instanceId] = ["Afreeca", user, nickname, format_bytes(file_size), format_duration(elapsed_time), output_filename]

##########################################################################
#################################  BIGO  #################################
##########################################################################

def bigo(instanceId, user):
  usernameList.insert(instanceId, ["Bigo", user, '', '', '', ''])

  verify = bigoVerify(user)

  if verify is True:
    m3u8Url = bigoPlaylist(user)
    siteId, nickname = bigoData(user)

  segment_urls = set()
  file_size = 0
  start_time = time.time()

  if os.path.exists('downloads/Bigo/' + user) == False:
    os.makedirs('downloads/Bigo/' + user)

  now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
  if platform.system() == 'Windows':
    now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

  output_filename = nickname + '-' + siteId + '-' + now + '-bigo.ts'
  output_path = 'downloads/Bigo/' + user + '/' + output_filename

  while True:
    base_url = m3u8Url.rsplit('/', 1)[0] + '/'

    res = requests.get(m3u8Url)
    playlist_content = res.text

    new_segment_lines = [
      line.strip() for line in playlist_content.splitlines() if line.endswith('.TS') or line.endswith('.ts')
    ]

    if '.ts' not in playlist_content.lower():
      try:
        usernameList[instanceId] = ["Bigo", user, "", 'Offline', 'Offline', 'Offline']

        if bigoVerify(user) is True:

          m3u8Url = bigoPlaylist(user)

          segment_urls = set()
          file_size = 0
          start_time = time.time()

          now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
          if platform.system() == 'Windows':
            now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())
          
          output_filename = nickname + '-' + siteId + '-' + now + '-bigo.ts'
          output_path = 'downloads/Bigo/' + user + '/' + output_filename
          
          continue
      except:
        usernameList[instanceId] = ["Bigo", user, "", 'ERR', 'ERR', 'ERR']

        if bigoVerify(user) is True:

          m3u8Url = bigoPlaylist(user)

          segment_urls = set()
          file_size = 0
          start_time = time.time()

          now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
          if platform.system() == 'Windows':
            now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())
          
          output_filename = nickname + '-' + siteId + '-' + now + '-bigo.ts'
          output_path = 'downloads/Bigo/' + user + '/' + output_filename
          
          continue

    with open(output_path, 'ab') as output_file:
        for new_segment_line in new_segment_lines:
          segment_url = urljoin(base_url, new_segment_line)
          if segment_url not in segment_urls:
            segment_urls.add(segment_url)
            try:
              res = requests.get(segment_url, timeout=15)
            except (requests.ReadTimeout, ConnectionError):
              continue
            file_size += len(res.content)
            elapsed_time = time.time() - start_time
            output_file.write(res.content)
            usernameList[instanceId] = ["Bigo", user, nickname, format_bytes(file_size), format_duration(elapsed_time), output_filename]

##########################################################################
################################# PANDA  #################################
##########################################################################
            
def panda(instanceId, user):
  usernameList.insert(instanceId, ["PandaTV", user, '', '', '', ''])

  verify = pandaVerify(user)

  if verify == "Err19":
    usernameList[instanceId] = ["PandaTV", user, '', 'Err19', 'Err19', 'INPUT COOKIES OF VERIFIED PANDA ACC']
    exit(1)

  if verify is True:
    m3u8Url = pandaPlaylist(user)

    segment_urls = set()
    file_size = 0
    start_time = time.time()

    if os.path.exists('downloads/PandaTV/' + user) == False:
      os.makedirs('downloads/PandaTV/' + user)

    now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
    if platform.system() == 'Windows':
      now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())

    output_filename = user + '-' + now + '-panda.ts'
    output_path = 'downloads/PandaTV/' + user + '/' + output_filename

    while True:
      base_url = m3u8Url.rsplit('/', 1)[0] + '/'

      res = requests.get(m3u8Url)
      playlist_content = res.text

      new_segment_lines = [
        line.strip() for line in playlist_content.splitlines() if line.endswith('.TS') or line.endswith('.ts')
      ]

      if '.ts' not in playlist_content.lower():
        try:
          usernameList[instanceId] = ["PandaTV", user, "", 'Offline', 'Offline', 'Offline']

          if pandaVerify(user) is True:
              
              m3u8Url = pandaPlaylist(user)

              segment_urls = set()
              file_size = 0
              start_time = time.time()

              now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
              if platform.system() == 'Windows':
                now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())
              
              output_filename = user + '-' + now + '-panda.ts'
              output_path = 'downloads/PandaTV/' + user + '/' + output_filename
              
              continue
        except:
          usernameList[instanceId] = ["PandaTV", user, "", 'ERR', 'RETRYING', 'CHECK COOKIES!']

          if pandaVerify(user) is True:
              
              m3u8Url = pandaPlaylist(user)

              segment_urls = set()
              file_size = 0
              start_time = time.time()

              now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())
              if platform.system() == 'Windows':
                now = time.strftime("%Y-%m-%d_%H-%M", time.localtime())
              
              output_filename = user + '-' + now + '-panda.ts'
              output_path = 'downloads/PandaTV/' + user + '/' + output_filename
              
              continue


      with open(output_path, 'ab') as output_file:
        for new_segment_line in new_segment_lines:
          segment_url = urljoin(base_url, new_segment_line)
          if segment_url not in segment_urls:
            segment_urls.add(segment_url)
            try:
              res = requests.get(segment_url, timeout=15)
            except (requests.ReadTimeout, ConnectionError):
              continue
            file_size += len(res.content)
            elapsed_time = time.time() - start_time
            output_file.write(res.content)
            usernameList[instanceId] = ["PandaTV", user, "", format_bytes(file_size), format_duration(elapsed_time), output_filename]