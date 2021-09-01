package main

import (
	"log"
	"os"
	"strings"

	"git.m2.nz/idanoo/nzcovidbot/internal/nzcovidbot"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting NZ Covid Bot")

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Fetch discord webhooks
	nzcovidbot.DiscordWebhooks = strings.Split(os.Getenv("DISCORD_WEBHOOKS"), ",")

	// Fetch slack webhook
	nzcovidbot.SlackWebhook = os.Getenv("SLACK_WEBHOOK")

	// Fetch twitter keys
	nzcovidbot.TwitterCreds = nzcovidbot.TwitterCredentials{
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	}

	// Git repo URL
	nzcovidbot.Repository = os.Getenv("SOURCE_REPO")

	// Boot up listeners / main loop
	nzcovidbot.Lesgoooo()
}
