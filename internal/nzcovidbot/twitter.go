package nzcovidbot

import (
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterCredentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

var TwitterCreds TwitterCredentials

func postToTwitter() {
	if TwitterCreds.AccessTokenSecret == "" {
		return
	}

	config := oauth1.NewConfig(TwitterCreds.ConsumerKey, TwitterCreds.ConsumerSecret)
	token := oauth1.NewToken(TwitterCreds.AccessToken, TwitterCreds.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	for _, row := range updatedLocations.Locations {
		_, _, err := client.Statuses.Update(row.TwitterData, nil)
		if err != nil {
			log.Print(err)
		}

		// Lets not ratelimit ourselves :upsidedownsmiley:
		time.Sleep(1 * time.Second)
	}
}
