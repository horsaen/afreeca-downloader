from plugins.afreeca.download import download, downloadVod
from plugins.afreeca.getPlaylist import getVodPlaylist, getVideoPlaylist, getStationNo
from plugins.afreeca.verify import verify

def main(fromStart, username, password):
  if fromStart is True:
    station_no = getStationNo(username, password)
    if station_no is False:
      station_no = input("Unable to get get stream id automatically, please input one manually:\n")
      vodUrl, filename = getVodPlaylist(username, station_no)
      downloadVod(vodUrl, filename, username)
  else:
    if verify(username):
      download(getVideoPlaylist(username, password), username)