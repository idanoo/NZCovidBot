# Changelog

## 1.2
- Change lastUpdated to pointer reference and use built in .After() functions
- Change date format to support multiple dates in single point of interest
- Update lastUpdated to max of previous event

## 1.1
- Split messages based on location
- Added systemd service
- Default to current time instead of 1970 (Unix timestamp 0)

## 1.0
- Reworked to use API
- Stores lastPoll / lastUpdated in lastUpdated.txt to keep track
- Only posts if PublishedDate > lastUpdated timestamp
- Removed twitter + git code
- Added more code comments