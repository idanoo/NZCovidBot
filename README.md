# NZCovidBot
Pull data from Github and parse new locations and fire them off to any configured endpoints

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
