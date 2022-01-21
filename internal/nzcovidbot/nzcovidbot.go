package nzcovidbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// Time of last succesful poll
var lastUpdated *time.Time

// Main func
func Lesgoooo() {
	// Set last updated poll time!
	lastUpdated = getLastPollTime()
	log.Printf("Using last updated time: %s", lastUpdated.Local())

	// Create chan to end timer
	endTicker := make(chan bool)

	// Timer to run every minute
	minuteTicker := time.NewTicker(time.Duration(60) * time.Second)

	// Initial poll check on load
	go getNewAPILocations()

	for {
		select {
		case <-endTicker:
			fmt.Println("Stopping background workers")
			return
		case <-minuteTicker.C:
			// Check for updates
			go getNewAPILocations()
		}
	}
}

// getLastPollTime - If run previously, get last TS, otherwise Now()
func getLastPollTime() *time.Time {
	// Set default of *now* if never run so we don't spam everything
	lastPoll := time.Now()

	// Load up last-polled date if set
	file, err := os.Open("lastUpdated.txt")
	if err == nil {
		b, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Unable to read lastUpdated.txt: %s", err)
			return &lastPoll
		}

		i, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			log.Printf("Unable to read lastUpdated.txt: %s", err)
			return &lastPoll
		}

		lastPoll = time.Unix(i, 0)
	}

	return &lastPoll
}

func postTheUpdates() {
	// Telegram
	go postToTelegram()

	// Slack
	go postToSlack()

	// Discord
	go postToDiscord()
}
