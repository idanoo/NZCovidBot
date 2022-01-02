package nzcovidbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

const API_ENDPOINT = "https://api.integration.covid19.health.nz/locations/v1/current-locations-of-interest"

var newLocations PostResponse

// Response from MoH API
type ApiResponse struct {
	Items []ApiItem `json:"items"`
}

// PostResponse - Above items ordered by location
type PostResponse struct {
	Items map[string][]ApiItem `json:"items"`
}

type ApiItem struct {
	EventID          string      `json:"eventId"`
	EventName        string      `json:"eventName"`
	StartDateTime    time.Time   `json:"startDateTime"`
	EndDateTime      time.Time   `json:"endDateTime"`
	PublicAdvice     string      `json:"publicAdvice"`
	VisibleInWebform bool        `json:"visibleInWebform"`
	PublishedAt      time.Time   `json:"publishedAt"`
	UpdatedAt        interface{} `json:"updatedAt"`
	ExposureType     string      `json:"exposureType"`
	Location         struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
		Suburb    string `json:"suburb"`
		City      string `json:"city"`
		Address   string `json:"address"`
	} `json:"location"`
}

// fetchAPILocations - Return struct of API response
func fetchAPILocations() (ApiResponse, error) {
	var apiResponse ApiResponse

	// Build HTTP Client and create request
	client := &http.Client{}
	req, err := http.NewRequest("GET", API_ENDPOINT, nil)
	if err != nil {
		return apiResponse, err
	}

	// Set user-agent info
	req.Header.Set("User-Agent", "NZCovidBot/1.0 (https://m2.nz)")

	// Fire off the request
	resp, err := client.Do(req)
	if err != nil {
		return apiResponse, err
	}
	defer resp.Body.Close()

	// Read body response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return apiResponse, err
	}

	// Unmarshal JSON into Go struct
	err = json.Unmarshal(body, &apiResponse)

	return apiResponse, err
}

// getNewAPILocations - Gets all locations and triggers posts
func getNewAPILocations() {
	// Set lastUpdate time at the start of the request so we don't miss any locations
	// posted right after we poll
	newPollTime := time.Now()

	// Pull latest data
	locations, err := fetchAPILocations()
	if err != nil {
		log.Printf("Error fetching API Locations %s", err)
		return
	}

	// Re-init our apiRepsonse so we don't hold onto old locations!
	newItems := make(map[string][]ApiItem, 0)

	// Iterate over the data and only find new locations
	for _, item := range locations.Items {
		if item.PublishedAt.Unix() > lastUpdated.Unix() {
			// Clone the item to put in our own lil slice
			copy := item
			newItems[item.Location.City] = append(newItems[item.Location.City], copy)
		}
	}

	// Make sure to clear out the previous list and append new data in a map based on location
	newLocations = PostResponse{}
	newLocations.Items = make(map[string][]ApiItem, 0)

	for mapKey, mapItems := range newItems {
		// Add location to our newLocations map
		newLocations.Items[mapKey] = mapItems

		// Order by StartDate
		sort.Slice(newLocations.Items[mapKey], func(i, j int) bool {
			return newLocations.Items[mapKey][i].StartDateTime.Before(newLocations.Items[mapKey][j].StartDateTime)
		})
	}

	// If new items, post it!
	if len(newLocations.Items) > 0 {
		postTheUpdates()
	}

	updateLastUpdated(newPollTime)
}

// updateLastUpdated - Creates/Updates lastUpdated.txt
func updateLastUpdated(newUpdated time.Time) {
	// Make sure to update the global var for next poll
	lastUpdated = newUpdated

	// Open file in truncate/append mode
	f, err := os.OpenFile("lastUpdated.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println(err)
		return
	}

	// Write data
	data := []byte(fmt.Sprintf("%d", lastUpdated.Unix()))
	_, err = f.Write(data)
	if err != nil {
		log.Println(err)
		return
	}

	// Close file so we can reopen next time
	if err := f.Close(); err != nil {
		log.Println(err)
	}
}

// getDateString - Returns Date + StartTime + EndTime
func (item ApiItem) getDateString() string {
	st := item.StartDateTime
	et := item.EndDateTime

	dateFormat := "Mon 2 Jan, 03:04PM"
	timeFormat := "03:04PM"

	return st.Local().Format(dateFormat) + " - " + et.Local().Format(timeFormat)
}