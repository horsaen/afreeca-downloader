# afreeca-downloader

A simple tool to download streams from afreeca and other sites

![tui screenshot](https://raw.githubusercontent.com/horsaen/imgstore/main/afreeca-downloader/tui.png)

Feature request or there's an issue the wiki doesn't cover? Join the Discord [here](https://discord.gg/yzNM7CPn4s)

## Supported Sites

The -mode flag supports the following sites

- afreeca
- bigo
- chzzk !!! REQUIRES YT-DLP !!!
- flex
- kick
- panda
- soop !!! REQUIRES YT-DLP !!!
- tiktok

Remember, some sites may require cookies to be able to download streams, read more below

## Installation and Usage

Have no clue what you're looking at? Don't worry! There's a simple guide [here](docs/getting-started.md) to get you started.

Pre-built binaries are available in the [releases](https://github.com/horsaen/afreeca-downloader/releases) tab.

To build from source:

Clone repo
```bash
git clone https://github.com/horsaen/afreeca-downloader.git
```

Install
```bash
go install
```

## Concurrent Downloads

![concurrent screenshot](https://raw.githubusercontent.com/horsaen/imgstore/main/afreeca-downloader/concurrent.png)

Supports all sites except:
- Chzzk
- Soop
- TikTok

To add users, open the users file @ `.afreeca-downloader/users`, and input users following the format

USERNAME, PLATFORM

```
user1, afreeca
user2, bigo
user3, flex
```

Then simply run the program using -concurrent

## Set Cookies

In order to function correctly, sometimes sites will require that you use cookies in order to validate network requests.

Found in your home folder @ `.afreeca-downloader/cookies`, you can input your corresponding cookies.

The cookies can be found in Developer tools > Storage > Cookies

!!WILL SOON SUPPORT AUTOCOOKIES!!

### Afreeca

rip afreeca, will update as needed, but will keep until it dies

- Copy only the VALUE of PdboxTicket into .afreeca-downloader/cookies/afreeca
- Should end up looking something like ``.A32.7bbT5``

### Panda

- Copy only the VALUE of sessKey into .afreeca-downloader/cookies/panda
- Should just look like a random string

Please note: using cookies on panda causes you to be kicked out of the current tab if you are grabbing the info for the same stream, there doesn't seem to be a fix for this (more testing needed)

### kick
This one requires a bit more cookies, please add the VALUES to .afreeca-downloader/cookies/kick in the following order:

1. __cf_bm
2. cf_clearance
3. kick_session

Mess up the order and it won't work

Please note: i don't actually know if you need them or not, so just add them if you encounter any issues

### FlexTV
- Copy only the VALUE of flx_oauth_access into .afreeca-downloader/cookies/flex
- Should just look like a random string

### Soop

- Copy only the VALUE of client-id into .afreeca-downloader/cookies/soop