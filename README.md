# afreeca-downloader

A simple tool to download streams from afreeca and other sites

Rewrite of [afreeca-downloader-old](https://github.com/horsaen/afreecatv-downloader-old), no longer using libraries such as streamlink, using a more reliable solution :DD

## Features
- DVR-like recording, neatly sorted by user in a download directory
- Downloading stream from start [experimental, but working]
- Multi-site support
- Batch downloading [not implemented]

## Installation and Usage

Clone repo
```bash
git clone https://github.com/horsaen/afreeca-downloader.git
```

Cd and use
```bash
cd afreeca-downloader && python3 main.py -h
```

## Set cookies

Sometimes cookies are needed to access certain data, they can be found in Developer tools > Storage > Cookies

Unless specified, cookies aren't needed

### afreeca
- Copy only the VALUE of PdboxTicket into [cookies](cookies)
- Should end up looking something like ``.A32.7bbT5``

### pandatv
- Copy only the VALUE of sessKey into [panda-cookies](plugins/pandatv/panda-cookies)
- Should just look like a random string

Please note: using cookies on panda causes you to be kicked out of the current tab, but not logged out, there does not seem to be a fix for this.A32.7bbT5

## Modes

the --mode flag supports the following arguements

- afreeca (default)
- panda
- bigo