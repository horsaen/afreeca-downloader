# Getting Started

## Easy Install

- Go to the [releases page](https://github.com/horsaen/afreeca-downloader/releases) and download right binary. If you are on Windows, download afreeca-downloader.exe.

- Simply run it and you're ready to download whatever !!

### In the case of cookies

Sometimes a given site will need cookies in order to properly make network requests. The standard workflow to obtain cookies for pretty much any site works like this:

1. Log into the site
2. Open up dev tools (F12)
3. Find cookies
- On Chome/chromium based browsers, this is found in Applications>Cookies
- Firefox browsers this is found in Storage>Cookies
4. Copy the value needed (listed in the README), for Afreeca it's PdboxTicket, it should look like .A32blablabla
5. Open the cookie file and paste your cookie. All the cookie files can be found in C:\Users\USERNAME\.afreeca-downloader/cookies, or for other OSs, wherever the home dir is.
6. Now interact with the TUI like normal, arrow keys + enter allow you to interact with it.
7. If any errors occur after the cookies are in use, it's likely they have expired, simply try the process again.
- If this hasn't worked, please make an issue on the Github or make a PR!