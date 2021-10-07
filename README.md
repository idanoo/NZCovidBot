# NZCovidBot
Pull data from Github and parse new locations and fire them off to any configured endpoints (Slack, Discord, Twitter(untested))

### About
After the twitterbot @nzcovidlocs shut down, I decided to try a different approach, instead of scraping MoH's website, lets parse the raw data!
https://github.com/minhealthnz/nz-covid-data/tree/main/locations-of-interest/august-2021    
It will clone ministry of healths git repo and poll every minute for updates to their raw CSV

## Config
Copy .env.example to .env and fill in the blanks.

### Run locally
```
    go run cmd/nzcovidbot/*.go
```

### Build
```
    go build -o nzcovidbot cmd/nzcovidbot/*.go
    ./nzcovidbot
```
