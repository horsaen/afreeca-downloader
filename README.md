# afreeca-downloader

A simple tool to download streams from afreeca and other sites

Rewrite of [afreeca-downloader-old](https://github.com/horsaen/afreecatv-downloader-old), no longer using libraries such as streamlink, using a more reliable solution :DD

## Features
- DVR-like recording, neatly sorted by user in a download directory
- Downloading stream from start
- Multi-site support
- Batch downloading [BETA]

## Installation and Usage

Clone repo
```bash
git clone https://github.com/sillysolutions/afreeca-downloader.git
```

Install deps
```bash
pip3 install -r requirements.txt
```

Cd and use
```bash
cd afreeca-downloader && python3 main.py -h
```

## Set cookies

Sometimes cookies are needed to access certain data, they can be found in Developer tools > Storage > Cookies

Unless specified, cookies aren't needed

### afreeca
- Copy only the VALUE of PdboxTicket into [cookies/afreeca](cookies/afreeca)
- Should end up looking something like ``.A32.7bbT5``

### pandatv
- Copy only the VALUE of sessKey into [cookies/panda](cookies/panda)
- Should just look like a random string

Please note: using cookies on panda causes you to be kicked out of the current tab if you are grabbing the info for the same stream, there doesn't seem to be a fix for this (more testing needed)

### kick
This one requires a bit more cookies, please add the VALUES to [cookies/kick](cookies/kick) in the following order:

1. __cf_bm
2. cf_clearance
3. kick_session

Mess up the order and it won't work

Please note: still heavily in beta, their api sucks and the reliability is horrid, don't expect it to work too much
New note: apparently it works a lot better if you're logged in, but needs more testing, keeping kick dl in the downloader for now

### FlexTV
- Copy only the VALUE of flx_oauth_access into [cookies/flex](cookies/flex)
- Should just look like a random string

## Modes

the --mode flag supports the following arguements

- afreeca (default)
- panda
- bigo
- kick
- flex
- tiktok

## Concurrent Downloads [BETA]

You can download multiple streams at once by using the --concurrent flag

Users to download from are specified in the [users](users) file, and follow the format of `username, platform`

For now the supported sites are

- afreeca
- bigo
- panda