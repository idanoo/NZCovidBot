package nzcovidbot

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
)

// Slack webhook URL
var SlackWebhooks []string

func postToSlack() {
	if len(SlackWebhooks) == 0 {
		return
	}

	postableData := getPostableSlackData()
	if len(postableData) == 0 {
		return
	}

	rawText := strings.Join(postableData, "\n")
	attachment1 := slack.Attachment{}
	attachment1.Text = &rawText

	payload := slack.Payload{
		Text:        "New Locations of Interest!",
		Username:    "NZCovidTracker",
		Channel:     "#covid-19-locations",
		IconUrl:     "https://www.skids.co.nz/wp-content/uploads/2020/08/download.png",
		Attachments: []slack.Attachment{attachment1},
	}

	for _, webhook := range SlackWebhooks {
		if webhook == "" {
			continue
		}

		err := slack.Send(webhook, "", payload)
		if len(err) > 0 {
			fmt.Print(err)
		}
	}
}
