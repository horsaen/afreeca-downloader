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

Set cookies
- Developer tools > Storage > Cookies
- Copy PdboxTicket into [cookies](cookies)
- Should end up looking something like ``PdboxTicket=.A32.``

## Supported sites
- [AfreecaTV](https://afreecatv.com/)
- [PandaTV](https://www.pandalive.co.kr/) [BETA, WORKING]
- [FlexTV](https://www.flextv.co.kr/) [WIP, NOT IMPLEMENTED]
- [Twitch](https://twitch.tv/) [WIP, NOT IMPLEMENTED]
- [Kick](https://kick.com/) [WIP, NOT IMPLEMENTED]
- [YouTube](https://youtube.com) [???]