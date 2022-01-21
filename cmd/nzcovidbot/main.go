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

	// Fetch slack webhooks
	nzcovidbot.SlackWebhooks = strings.Split(os.Getenv("SLACK_WEBHOOKS"), ",")

	// Fetch slack webhooks
	nzcovidbot.TelegramTokens = strings.Split(os.Getenv("TELEGRAM_BOT_TOKENS"), ",")

	// Boot up listeners / main loop
	nzcovidbot.Lesgoooo()
}
