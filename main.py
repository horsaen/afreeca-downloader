from utils.download import download
from utils.downloadVod import downloadVod
from utils.getPlaylist import getVodPlaylist, getVideoPlaylist, getStationNo
from utils.verify import verify
import argparse
import os

from plugins.pandatv.panda import main as pandaMain
from plugins.afreeca_m3u8.download_m3u8 import download_m3u8

# qualities::::
# original
# hd
# sd

def flagsInit():
    parser = argparse.ArgumentParser(description="Afreeca TV Downloader :D")

    parser.add_argument('-u', '--username', default=False, help='Streamer username')
    parser.add_argument('-p', '--password', default='', help='Stream password [EXPERIMENTAL]')
    parser.add_argument('--from-start', default=False, action='store_true', help='Download from the start of the stream [EXPERIMENTAL]')
    parser.add_argument('--panda', default=False, action='store_true', help='Download video from PandaTV [NOT IMPLEMENTED]')
    parser.add_argument('--batch', default=False, action='store_true', help='Download multiple streams from a text file [NOT IMPLEMENTED]')
    parser.add_argument('--playlist', default=False, help="Download from an afreeca m3u8 playlist, anything other than smil:vod has not been tested")

    args = parser.parse_args()
    return args

def main(username, pwd, args):
    if args.panda is True:
        pandaMain(username)
    if args.playlist is not False:
        download_m3u8(args.playlist)
    if args.from_start is True:
        station_no = getStationNo(username, pwd)
        if station_no is False:
            station_no = input("Unable to get get stream id automatically, please input one manually:\n")
        vodUrl, filename = getVodPlaylist(username, station_no)
        downloadVod(vodUrl, filename, username)
    if verify(username):
        download(getVideoPlaylist(username, pwd), username)

if __name__ == '__main__':
    args = flagsInit()

    if os.path.exists('downloads') is False:
        os.makedirs('downloads')

    username = args.username
    pwd = args.password

    if username is False and args.playlist is False:
        username = input('Enter username:\n')

    # # change to switch case later
    # if panda is True:
        # pandaMain(username)
    # else:
    main(username, pwd, args)