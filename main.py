import argparse
import os

from plugins.afreeca import mainProcess as afreeca
from plugins.bigo import mainProcess as bigo
from plugins.concurrent import main as concurrent
from plugins.flex import mainProcess as flex
from plugins.kick import downloadStream as kick
from plugins.panda import mainProcess as panda
from plugins.tiktok import mainProcess as tiktok

# qualities::::
# original
# hd
# sd

def flagsInit():
  parser = argparse.ArgumentParser(description="Afreeca TV Downloader :D")

  parser.add_argument('-u', '--username', default=False, help='Streamer username')
  parser.add_argument('-p', '--password', default='', help='Stream password [EXPERIMENTAL]')
  parser.add_argument('-m', '--mode', default='afreeca', help='Select site, supported sites are found on the github readme')
  parser.add_argument('--from-start', default=False, action='store_true', help='Download from the start of the stream [EXPERIMENTAL], only working for afreeca')
  parser.add_argument('--playlist', default=False, action='store_true', help="Download from an afreeca m3u8 playlist, anything other than smil:vod has not been tested")
  parser.add_argument('--concurrent', default=False, action='store_true', help='Download multiple streams concurrently [EXPERIMENTAL]')

  args = parser.parse_args()
  return args

def main(args):
  username = args.username
  pwd = args.password
  if username is False and args.concurrent is False and args.playlist is False:
    username = input('Enter username:\n')
  if args.concurrent is True:
    concurrent()
  elif args.mode == 'afreeca':
    mode = "live"
    if args.from_start is True: mode = "fromStart"
    elif args.playlist is True: mode = "playlist"
    afreeca(mode, username, pwd)
  elif args.mode == 'bigo':
    bigo(username)
  elif args.mode == 'flex':
    flex(username)
  elif args.mode == 'kick':
    kick(username)
  elif args.mode == 'panda':
    panda(username)
  elif args.mode == 'tiktok':
    tiktok(username)
  # elif args.mode == 'twitch':
    # twitchMain(username)
  else:
    print('Invalid mode')

if __name__ == '__main__':
  args = flagsInit()

  if os.path.exists('downloads') is False:
    os.makedirs('downloads')

  main(args)