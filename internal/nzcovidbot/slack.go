package nzcovidbot

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
)

// Slack webhook URL
var SlackWebhook string

func postToSlack() {
	if SlackWebhook == "" {
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
		Channel:     "#covid-19",
		IconUrl:     "https://www.skids.co.nz/wp-content/uploads/2020/08/download.png",
		Attachments: []slack.Attachment{attachment1},
	}

	err := slack.Send(SlackWebhook, "", payload)
	if len(err) > 0 {
		fmt.Printf("error: %s\n", err)
	}
}
