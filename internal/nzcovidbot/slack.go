package nzcovidbot

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
)

// Slack webhook#channel
var SlackWebhooks []string

// Send slack request
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

// Adds new rows to a slice for slack
func getPostableSlackData() []slack.Attachment {
	rows := make([]slack.Attachment, 0)
	if len(newLocations.Items) == 0 {
		return rows
	}

	for _, item := range newLocations.Items {
		rows = append(rows, getSlackRow(item))
	}

	return rows
}

// getSlackRow - Get slack attachment row
func getSlackRow(item ApiItem) slack.Attachment {
	url := getMapsLinkFromAddress(item.EventName, item.Location.Address)
	dateRange := item.getDateString()

	attachment := slack.Attachment{
		Title:     &item.EventName,
		TitleLink: &url,
		Text:      &dateRange,
	}

	return attachment
}

// getMapsLinkFromAddress hyperlink gmaps
func getMapsLinkFromAddress(name string, address string) string {
	return fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%s", url.QueryEscape(name+", "+address))
}
