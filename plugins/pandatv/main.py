import os
import requests

from plugins.pandatv.download import download
from plugins.pandatv.verify import verify

def main(username):
  if os.path.exists('downloads/PandaTV') is False:
    os.makedirs('downloads/PandaTV')

  download(verify(username), username)