import streamlink

def manual_quality(streamUrl):
    streams = streamlink.streams(streamUrl)
    qualityList = list(streams.keys())
    return qualityList