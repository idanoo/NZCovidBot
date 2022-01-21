package nzcovidbot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Slice of discord webhooks
var TelegramTokens []string

// Max locations per telegram post
var TelegramMaxPerPost = 100

func postToTelegram() {
	postableTelegramData := getPostableTelegramData()
	if len(postableTelegramData) == 0 {
		return
	}

	for _, telegramToken := range TelegramTokens {
		if telegramToken != "" {
			tokenParts := strings.Split(telegramToken, "!")
			len := len(tokenParts)
			if len != 2 {
				log.Println("Telegram token error, channel ID or token invalid")
				continue
			}

			bot, err := tgbotapi.NewBotAPI(tokenParts[len-2])
			if err != nil {
				log.Print(err)
				continue
			}

			// Build message and send for each location
			for location, postableData := range postableTelegramData {
				for _, post := range postableData {
					chanId, err := strconv.ParseInt(tokenParts[len-1], 10, 64)
					if err != nil {
						continue
					}

					// Create msg
					msg := tgbotapi.NewMessage(chanId, "<b><u>"+telegramEscape(location)+"</u></b>\n\n"+post)

					// Decided to use HTML here as Markdown has way too many restricted chars :<
					msg.ParseMode = tgbotapi.ModeHTML

					// SEND IT
					_, err = bot.Send(msg)
					if err != nil {
						log.Print(err)
					}

					time.Sleep(500 * time.Millisecond)
				}
			}
		}
	}
}

// getPostableTelegramData - Returns slices containing 20~ locations each
// to send as separate messages. map[location][]locationsofinterest
func getPostableTelegramData() map[string][]string {
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
			rows = append(rows, getTelegramRow(item))

			// Make sure to create a new slice if we have >100 to send as a different message
			if len(rows) > TelegramMaxPerPost {
				groups[location] = append(groups[location], strings.Join(rows, "\n\n"))
				rows = make([]string, 0)
			}
		}

		// If we have less than 100, append any more before next location
		if len(rows) > 0 {
			groups[location] = append(groups[location], strings.Join(rows, "\n\n"))
		}
	}

	return groups
}

// getTelegramRow Format the string to a tidy string for the interwebs
func getTelegramRow(item ApiItem) string {
	return fmt.Sprintf("<b>%s</b>\n%s\n<i>%s</i>",
		telegramEscape(item.EventName), telegramEscape(item.Location.Address), telegramEscape(item.getDateString()))
}

func telegramEscape(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				s,
				"&",
				"&amp;",
			),
			">",
			"&gt;",
		),
		"<",
		"&lt;",
	)
}
