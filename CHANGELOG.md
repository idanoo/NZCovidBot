# Changelog

## 1.1
- Split messages based on location
- Added systemd service

## 1.0
- Reworked to use API
- Stores lastPoll / lastUpdated in lastUpdated.txt to keep track
- Only posts if PublishedDate > lastUpdated timestamp
- Removed twitter + git code
- Added more code comments