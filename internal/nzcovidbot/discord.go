package nzcovidbot

import (
	"log"
	"strings"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
)

var (
	// Slice of discord webhooks
	DiscordWebhooks []string
)

func postToDiscord(webhookString string, msg string) {
	tokenParts := strings.Split(webhookString, "/")
	len := len(tokenParts)
	webhook, err := disgohook.NewWebhookClientByToken(nil, nil, tokenParts[len-2]+"/"+tokenParts[len-1])
	if err != nil {
		log.Print(err)
		return
	}

	_, err = webhook.SendEmbeds(api.NewEmbedBuilder().
		SetDescription(msg).
		Build(),
	)
	if err != nil {
		log.Print(err)
		return
	}

	if err != nil {
		log.Print(err)
	}
}
