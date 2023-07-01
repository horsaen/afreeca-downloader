# migrate afreeca logic to here !!
import streamlink
import time

# import our own modules
from verify import verify
from utils.formatBytes import format_bytes
from utils.formatDuration import format_duration

def downloadStream(username, streamUrl, qualityName, outputFlag):

    session = streamlink.Streamlink()
    session.set_option("stream-segment-threads", 10)
    session.set_plugin_option("afreeca", "username", "")
    session.set_plugin_option("afreeca", "password", "")

    streams = session.streams(streamUrl)
    stream = streams[qualityName]
    
    print("Downloading stream using " + qualityName + " quality")

    now = time.strftime("%Y-%m-%d_%H:%M", time.localtime())

    outputPath: str = username + "-" + now + "-afreeca.ts" 

    # if outputFlag is not False:
    #     outputPath = outputFlag

    with stream.open() as fd:
        with open(outputPath, 'wb') as output:
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
                    print("\r" + f"Downloading to {outputPath} || {format_duration(elapsed_time)} @ {format_bytes(file_size)}", end="", flush=True)
                except Exception as e:
                    print("An error occured: " + str(e))
                    if verify(username, urlStruct) is not None:
                        downloadStream(username, streamUrl, qualityName, output)