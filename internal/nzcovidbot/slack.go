package nzcovidbot

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
)

// Slack webhook#channel
var SlackWebhooks []string

func postToSlack() {
	if len(SlackWebhooks) == 0 {
		return
	}

	attachmentData := getPostableSlackData()
	if len(attachmentData) == 0 {
		return
	}

	for _, webhook := range SlackWebhooks {
		if webhook == "" {
			continue
		}

		parts := strings.Split(webhook, "!")
		payload := slack.Payload{
			Text:        "New Locations of Interest!",
			Username:    "NZCovidTracker",
			Channel:     "#" + parts[1],
			IconUrl:     "https://www.skids.co.nz/wp-content/uploads/2020/08/download.png",
			Attachments: attachmentData,
		}

		err := slack.Send(parts[0], "", payload)
		if len(err) > 0 {
			fmt.Printf("Wehbook: %s\nError: %s", webhook, err)
		}
	}
}
