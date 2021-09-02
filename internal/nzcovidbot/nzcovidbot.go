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
	// Lets reshuffle our structured data a bit (Exposure Date ASC)
	orderRowDataByDate()

	// Twitter
	go postToTwitter()

	// Slack
	go postToSlack()

	// Discord
	postableDiscordData := getPostableDiscordData()
	if len(postableDiscordData) == 0 {
		return
	}

	for _, discordWebhook := range DiscordWebhooks {
		for _, postableData := range postableDiscordData {
			go postToDiscord(discordWebhook, postableData)
			time.Sleep(1 * time.Second)
		}
	}

	// Clear out posted data!
	updatedLocations = UpdatedLocations{}
}
