package nzcovidbot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
)

// Slice of discord webhooks
var DiscordWebhooks []string

func postToDiscord() {
	postableDiscordData := getPostableDiscordData()
	if len(postableDiscordData) == 0 {
		return
	}

	for _, discordWebhook := range DiscordWebhooks {
		for _, postableData := range postableDiscordData {
			if discordWebhook != "" {
				tokenParts := strings.Split(discordWebhook, "/")
				len := len(tokenParts)

				// Build discord request
				webhook, err := disgohook.NewWebhookClientByToken(nil, nil, tokenParts[len-2]+"/"+tokenParts[len-1])
				if err != nil {
					log.Print(err)
					return
				}

				// Send discord message
				_, err = webhook.SendEmbeds(api.NewEmbedBuilder().
					SetDescription(postableData).
					Build(),
				)

				if err != nil {
					log.Print(err)
				}

				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

// getPostableDiscordData - Returns slices containing 20~ locations each
// to send as separate messages
func getPostableDiscordData() []string {
	// Create our slice of groups
	groups := make([]string, 0)
	if len(newLocations.Items) == 0 {
		return groups
	}

	rows := make([]string, 0)
	for _, item := range newLocations.Items {
		rows = append(rows, getDiscordRow(item))

		if len(rows) > 20 {
			groups = append(groups, strings.Join(rows, "\n"))
			rows = make([]string, 0)
		}
	}

	return append(groups, strings.Join(rows, "\n"))
}

// formatCsvDiscordRow Format the string to a tidy string for the interwebs
func getDiscordRow(item ApiItem) string {
	return fmt.Sprintf("**%s** %s on _%s_",
		item.EventName, item.Location.Address, item.getDateString())

}
