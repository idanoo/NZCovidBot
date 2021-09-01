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
	// Discord
	postableDiscordData := getPostableDiscordData()

	if postableDiscordData == "" {
		return
	}

	for _, discordWebhook := range DiscordWebhooks {
		postToDiscord(discordWebhook, postableDiscordData)
	}

	// Twitter
	// postableTwitterData := getPostableTwitterData()
	// if postableTwitterData == "" {
	// 	return
	// }

	// for _, discordWebhook := range DiscordWebhooks {
	// 	postToTwitter(discordWebhook, postableTwitterData)
	// }

	// Clear out posted data!
	updatedLocations = UpdatedLocations{}
}
