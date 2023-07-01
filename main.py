import simple_term_menu
import streamlink
import requests
import os
import time
import argparse

# import our own modules
from verify import verify
from utils.formatBytes import format_bytes
from utils.formatDuration import format_duration
from utils.manual_quality import manual_quality

def flagsInit():
    parser = argparse.ArgumentParser(description="Afreeca TV Downloader :D")

    parser.add_argument("-m", "--manual", default=False, action="store_true", help="Manually select download quality")
    parser.add_argument("-o", "--output", default=False, help="File output, defaults to ")
    parser.add_argument("-su", "--susername", default=False, help="Streamer username")
    parser.add_argument("-u", "--username", default=False, help="Afreeca Username")
    parser.add_argument("-p", "--password", default=False, help="Afreeca Password")

    args = parser.parse_args()

    return args

def downloadStream(username, streamUrl, qualityName, outputFlag):

    session = streamlink.Streamlink()
    session.set_option("stream-segment-threads", 10)
    session.set_plugin_option("afreeca", "username", "")
    session.set_plugin_option("afreeca", "password", "")

    streams = session.streams(streamUrl)
    stream = streams[qualityName]
    
    print("Downloading stream using " + qualityName + " quality")

    now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())

    if os.path.exists("downloads/" + username) is False:
        os.makedirs("downloads/" + username)

    output_path: str = username + "-" + now + "-afreeca.ts"
    output_path_absolute: str = "downloads/" + username + "/" + username + "-" + now + "-afreeca.ts" 

    # if outputFlag is not False:
    #     outputPath = outputFlag

    with stream.open() as fd:
        with open(output_path_absolute, 'wb') as output:
            file_size = 0
            start_time = time.time()
            while True:
                try:
                    data = fd.read(1024)
                    if not data:
                        break
                    output.write(data)
                    file_size += len(data)
                    elapsed_time = time.time() - start_time
                    # ugly but works
                    print("\r" + f"Downloading to {output_path} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}             \x1b[?25l", end="", flush=True)
                    # loop !
                except Exception as e:
                    fd.close()
                    print("An error occured: " + str(e))
                    if verify(username, urlStruct) is not None:
                        downloadStream(username, streamUrl, qualityName, output)

if __name__ == '__main__':
    args = flagsInit()

    # check if dl folder exists :D
    if os.path.exists("downloads") is False:
        os.makedirs("downloads")

    username = args.username
    
    if username is False:
        username = input("Enter streamer username:\n")

    urlStruct: str = 'https://bjapi.afreecatv.com/api/' + username + '/station'

    streamUrl = verify(username, urlStruct)

    output = args.output

    print(streamUrl)

    qualityName: str = 'best'

    if args.manual is True:
        quality = manual_quality(streamUrl)
        menu = simple_term_menu.TerminalMenu(quality)
        menuIndex = menu.show()
        qualityName = quality[menuIndex]

    downloadStream(username, streamUrl, qualityName, output)