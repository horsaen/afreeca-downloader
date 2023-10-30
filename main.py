from utils.download import download
from utils.downloadVod import downloadVod
from utils.getPlaylist import getVodPlaylist, getVideoPlaylist, getStationNo
from utils.verify import verify
import argparse
import os

from plugins.pandatv.main import main as pandaMain
from plugins.pandatv.verify import checkUser as pandaVerify
from plugins.bigo.main import main as bigoLive
from plugins.kick.main import main as kickMain
from plugins.tt.main import main as ttMain
from plugins.twitch.main import main as twitchMain

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
    # parser.add_argument('--batch', default=False, action='store_true', help='Download multiple streams from a text file [NOT IMPLEMENTED]')
    parser.add_argument('--playlist', default=False, help="Download from an afreeca m3u8 playlist, anything other than smil:vod has not been tested")

    args = parser.parse_args()
    return args

def main(username, pwd, args):
    if args.mode == 'afreeca':
        if args.from_start is True:
            station_no = getStationNo(username, pwd)
            if station_no is False:
                station_no = input("Unable to get get stream id automatically, please input one manually:\n")
            vodUrl, filename = getVodPlaylist(username, station_no)
            downloadVod(vodUrl, filename, username)
        else:
            if verify(username):
                download(getVideoPlaylist(username, pwd), username)
    elif args.mode == 'panda':
        if pandaVerify(username):
            pandaMain(username)
    elif args.mode == 'bigo':
        bigoLive(username)
    elif args.mode == 'kick':
        kickMain(username)
    elif args.mode == 'tiktok':
        ttMain(username)
    elif args.mode == 'twitch':
        twitchMain(username)

if __name__ == '__main__':
    args = flagsInit()

    if os.path.exists('downloads') is False:
        os.makedirs('downloads')

    username = args.username
    pwd = args.password

    if username is False and args.playlist is False:
        username = input('Enter username:\n')

    main(username, pwd, args)