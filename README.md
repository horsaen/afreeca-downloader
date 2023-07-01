# afreeca-downloader

A simple tool for the automatic download of afreeca (and other sites) 

## Installation and Usage

Clone repo
```bash
git clone https://github.com/horsaen/afreeca-downloader
```

Cd into directory and install dependencies

```bash
cd afreeca-downloader
pip3 install -r REQUIREMENTS.txt
```

Run

```bash
python3 main.py -h
```

## Todo
- [ ] Allow for custom output directory
- [ ] Remove unneeded error warnings, particularly the streamlink warning about missing segments or something
- [ ] Take inspo from streamlink themselves and make this work for multiple sites (depsite the name, so funny) because that's cool
- [ ] Add to pip or whatever that is
- [ ] Multithreading

## FFMPEG workflow
Here's just a simple guide on changing the direct output to something a bit more storage friendly

Ignore .ts segment errors:

```bash
ffmpeg -i <INPUT>.ts -map 0 -ignore_unkown -c copy <OUTPUT>.ts
```

### Matroska

```bash
ffmpeg -i <INPUT>.ts -map 0 -c copy <OUTPUT>.mkv
```

### MP4

```bash
ffmpeg -i <INPUT>.ts -map 0 -c copy <OUTPUT>.mp4
```

Re-encode the video to H.264 and stream copy the audio:

```bash
ffmpeg -i input.ts -c:v libx264 -c:a copy output.mp4
```

Re-encode both video and audio:

```bash
ffmpeg -i input.ts -c:v libx264 -c:a aac output.mp4
```

Lossless H.264 example:

```bash
ffmpeg -i input.ts -c:v libx264 -crf 0 -c:a copy output.mp4
```

Lossless files will be huge ^^

### Lossless H.265

```bash
ffmpeg -i <INPUT>.ts -c:v libx265 -x265-params "lossless=1" <OUTPUT>.mkv
```