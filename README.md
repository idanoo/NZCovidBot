# NZCovidBot
Pulls data from Ministry of Health API and parse into Discord and Slack webhooks.

### About
After the twitterbot @nzcovidlocs shut down, I decided to try a different approach, instead of scraping MoH's website, we originally parsed the raw CSV data.
Since then the NZ Ministry of Health have released an API containing this data now. We are now using this https://api.integration.covid19.health.nz/locations/v1/current-locations-of-interest 

## Config
Copy .env.example to .env and fill in the webhook URLs

### Run locally
```
    go run cmd/nzcovidbot/*.go
```

### Build
```
    go build -o nzcovidbot cmd/nzcovidbot/*.go
    sudo cp nzcovidbot.service /etc/systemd/system/nzcovidbot.service
    # Update user + location of repo in systemd file
    sudo systemctl daemon-reload && systemctl enable --now nzcovidbot.service
```

### Screenshots

#### Discord
![DiscordExample](https://gitlab.com/idanoo/NZCovidBot/-/raw/master/discordexample.png)

#### Slack
![SlackExample](https://gitlab.com/idanoo/NZCovidBot/-/raw/master/slackexample.png)