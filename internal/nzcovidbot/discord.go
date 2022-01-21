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
		if discordWebhook != "" {
			// Build discord request
			tokenParts := strings.Split(discordWebhook, "/")
			len := len(tokenParts)
			webhook, err := disgohook.NewWebhookClientByToken(nil, nil, tokenParts[len-2]+"/"+tokenParts[len-1])
			if err != nil {
				log.Print(err)
				continue
			}

			// Build message and send for each location
			for location, postableData := range postableDiscordData {
				for _, post := range postableData {

					// Send discord message
					_, err = webhook.SendEmbeds(api.NewEmbedBuilder().
						SetTitle("*" + location + "*").
						SetDescription(post).
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
}

// getPostableDiscordData - Returns slices containing 20~ locations each
// to send as separate messages. map[location][]locationsofinterest
func getPostableDiscordData() map[string][]string {
	// Create our return map
	groups := make(map[string][]string, 0)

	// If no locations, lets return empty map
	if len(newLocations.Items) == 0 {
		return groups
	}

	for location, items := range newLocations.Items {
		// Create out output buffer per location
		rows := make([]string, 0)

		// Foreach item, create the output text based off the item
		for _, item := range items {
			rows = append(rows, getDiscordRow(item))

			// Make sure to create a new slice if we have >20 to send as a different message
			if len(rows) > 20 {
				groups[location] = append(groups[location], strings.Join(rows, "\n"))
				rows = make([]string, 0)
			}
		}

		// If we have less than 20, append any more before next location
		if len(rows) > 0 {
			groups[location] = append(groups[location], strings.Join(rows, "\n"))
		}
	}

	return groups
}

// getDiscordRow Format the string to a tidy string for the interwebs
func getDiscordRow(item ApiItem) string {
	return fmt.Sprintf("**%s** %s on _%s_",
		item.EventName, item.Location.Address, item.getDateString())
}
