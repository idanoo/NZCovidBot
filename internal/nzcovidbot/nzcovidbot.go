package nzcovidbot

import (
	"fmt"
	"time"
)

var Repository string

func Lesgoooo() {
	// Setup repo stuff
	loadRepo(Repository)

	// Create chan to end timer
	endTicker := make(chan bool)

	// Timer to run every minute
	minuteTicker := time.NewTicker(time.Duration(60) * time.Second)

	// Initial check on load
	go checkForUpdates()

	for {
		select {
		case <-endTicker:
			fmt.Println("Stopping background workers")
			return
		case <-minuteTicker.C:
			// Check for updates
			go checkForUpdates()
		}
	}

}

func postTheUpdates() {
	// Twitter
	go postToTwitter()

	// Discord
	postableDiscordData := getPostableDiscordData()
	if postableDiscordData == "" {
		return
	}

	// Not using go routines so we don't get rate limited
	for _, discordWebhook := range DiscordWebhooks {
		postToDiscord(discordWebhook, postableDiscordData)
	}

	// Clear out posted data!
	updatedLocations = UpdatedLocations{}
}
