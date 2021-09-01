package nzcovidbot

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"
	"time"
)

// Slice of updated located
type UpdatedLocations struct {
	Locations []UpdatedRow
}

// Updated data
type UpdatedRow struct {
	ChangeType  string // ADDED, REMOVED, MODIFIED
	DiscordData string // Formatted Row data
	TwitterData string // Formatted Row data
	SlackData   string // Formatted Row data
}

// Struct of updated locations
var updatedLocations UpdatedLocations

// parseCsvRow Build into struct for output later
func parseCsvRow(changeType string, data string) {
	newRow := UpdatedRow{
		ChangeType:  changeType,
		DiscordData: formatCsvDiscordRow(data),
		TwitterData: formatCsvTwitterRow(data),
		SlackData:   formatCsvSlackRow(data),
	}

	updatedLocations.Locations = append(updatedLocations.Locations, newRow)
}

// formatCsvDiscordRow Format the string to a tidy string for the interwebs
func formatCsvDiscordRow(data string) string {
	c := parseRawRowData(data)
	return fmt.Sprintf("**%s** %s on _%s_ - _%s_", c[2], c[3], c[0], c[1])
}

// formatCsvTwitterRow Format the string to a tidy string for the interwebs
func formatCsvTwitterRow(data string) string {
	c := parseRawRowData(data)
	return fmt.Sprintf("New Location: %s\n%s\n%s - %s\n#NZCovidTracker #NZCovid", c[2], c[3], c[0], c[1])
}

// formatCsvSlackRow Format the string to a tidy string for the interwebs
func formatCsvSlackRow(data string) string {
	c := parseRawRowData(data)
	return fmt.Sprintf("*%s* %s on _%s_ - _%s_", c[2], c[3], c[0], c[1])
}

// Returns []string of parsed
func parseRawRowData(data string) []string {
	output := make([]string, 0)

	r := csv.NewReader(strings.NewReader(data))
	r.Comma = ','
	fields, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return output
	}

	c := make([]string, 0)
	c = append(c, fields...)

	starttime := c[4]
	st, err := time.Parse("02/01/2006, 3:04 pm", starttime)
	if err != nil {
		log.Print(err)
	} else {
		starttime = st.Format("Monday 2 Jan, 3:04PM")
	}

	endtime := c[5]
	et, err := time.Parse("02/01/2006, 3:04 pm", endtime)
	if err != nil {
		log.Print(err)
		endtime = strings.Split(c[5], ", ")[1]
	} else {
		endtime = et.Format("3:04PM")
	}

	return append(output, starttime, endtime, c[1], c[2])
}

func getPostableDiscordData() string {
	if len(updatedLocations.Locations) == 0 {
		return ""
	}

	rows := make([]string, 0)
	for _, location := range updatedLocations.Locations {
		if location.ChangeType == "REMOVED" {
			rows = append(rows, fmt.Sprintf("REMOVED: %s", location.DiscordData))
		} else {
			rows = append(rows, location.DiscordData)
		}
	}

	return "\\`\\`\\`\n" + strings.Join(rows, "\n") + "\n\\`\\`\\`"
}

func getPostableSlackData() []string {
	rows := make([]string, 0)
	if len(updatedLocations.Locations) == 0 {
		return rows
	}

	for _, location := range updatedLocations.Locations {
		if location.ChangeType == "REMOVED" {
			rows = append(rows, fmt.Sprintf("REMOVED: %s", location.SlackData))
		} else {
			rows = append(rows, location.SlackData)
		}
	}

	return rows
}
