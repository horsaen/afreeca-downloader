import requests
from plugins.tt.getPlaylist import getStreamPlaylist, getRoomId
from plugins.tt.verify import verify
from plugins.tt.download import downloadPlaylist

def main(username):
  room_id = verify(username)
  playlist = getStreamPlaylist(room_id)
  downloadPlaylist(playlist, username)