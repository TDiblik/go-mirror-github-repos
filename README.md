## About

A program to automatically mirror GitHub repositories locally.
It finds all of the repositories you're part of, using the `https://api.github.com/user/repos` endpoint.
Then it goes and either clones them, or forcefully updates every local branch.
It doesn't have any sort of functionality to run automatically, so you have to setup crontab with it.

## Setup

- Create a new `Personal access token (classic)`
  - Inside `Settings` -> `Developer settings` -> `Personal access tokens` -> `Tokens (classic)` -> `Generate new token (classic)`
  - You don't have to give it any permissions
  - After it expires, you have to do this process again, personally, I set the expiration to `No expiration`, since I want to `set it & forget it`
- Create `.env`, you can use `.env.example` as a template
  - `GH_TOKEN` -> Your GitHub token you created
  - `MIRROR_PATH` -> A directory path where the program will mirror your GH. **THIS PATH MUST EXISTS**
  - `EXCLUDED_ORGS` -> A list of orgs/usernames to ignore (eg. `TDiblik, MyUsername, MyOrg, MyWorkOrg,..`)
  - `EXCLUDED_REPOSITORIES` -> A list of specific repositories to ignore (eg. `MyUsername/MyRepo, TDiblik/main-gate-alpr, TDiblik/KeyXpert,..`)
- Build the program and copy it onto your server (with the `.env`)
  - `GOOS=linux GOARCH=amd64 go build -o build/go-mirror-github-repos`
  - scp with `.env`: `scp -P {PORT} build/go-mirror-github-repos .env {USERNAME}@{IP}:~/go-mirror-github-repos/`
  - scp without `.env`: `scp -P {PORT} build/go-mirror-github-repos {USERNAME}@{IP}:~/go-mirror-github-repos/`
  - `cd ~/go-mirror-github-repos/; chmod +x go-mirror-github-repos;`
- Create the `MIRROR_PATH`
  - This is mine: `mkdir ~/go-mirror-github-repos/mirror/`
- Setup a cron job to run the program at a desired time
  - This is mine: `0 1 * * * cd /home/{USERNAME}/go-mirror-github-repos && PATH=/usr/local/bin:/usr/bin:/bin /home/{USERNAME}/go-mirror-github-repos/go-mirror-github-repos > /home/{USERNAME}/go-mirror-github-repos/log.txt 2>&1`
