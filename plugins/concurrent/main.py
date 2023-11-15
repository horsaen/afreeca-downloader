import time, os
from threading import Thread
from tabulate import tabulate

from plugins.concurrent.downloadStream import downloadStream, usernameList

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
      
      thread = Thread(target=downloadStream, args=(instanceId, name, platform))
      thread.start()
      threads.append(thread)
  
  while True:
    try:
      head = ["Site", "User", "Nick", "Size", "Duration", "Path"]
      os.system('cls' if os.name == 'nt' else 'clear')
      # print("\r" + str(usernameList) + '\x1b[?25l', end='', flush=True)
      # print(usernameList)
      print(tabulate(usernameList, headers=head, tablefmt='grid') + '\x1b[?25l')
      time.sleep(2)
    except KeyboardInterrupt:
      exit()