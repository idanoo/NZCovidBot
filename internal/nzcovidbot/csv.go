package nzcovidbot

import (
	"encoding/csv"
	"fmt"
	"strings"
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
}

// Struct of updated locations
var updatedLocations UpdatedLocations

// parseCsvRow Build into struct for output later
func parseCsvRow(changeType string, data string) {
	newRow := UpdatedRow{
		ChangeType:  changeType,
		DiscordData: formatCsvDiscordRow(data),
		TwitterData: formatCsvTwitterRow(data),
	}

	updatedLocations.Locations = append(updatedLocations.Locations, newRow)
}

// formatCsvDiscordRow Format the string to a tidy string for the interwebs
func formatCsvDiscordRow(data string) string {
	// Split string
	r := csv.NewReader(strings.NewReader(data))
	r.Comma = ','
	fields, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	c := make([]string, 0)
	c = append(c, fields...)
	endtime := strings.Split(c[5], ", ")
	return fmt.Sprintf("**%s** %s on _%s_ - _%s_", c[1], c[2], c[4], endtime[1])
}

// formatCsvTwitterRow Format the string to a tidy string for the interwebs
func formatCsvTwitterRow(data string) string {
	// Split string
	r := csv.NewReader(strings.NewReader(data))
	r.Comma = ','
	fields, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	c := make([]string, 0)
	c = append(c, fields...)
	endtime := strings.Split(c[5], ", ")
	return fmt.Sprintf("**%s** - %s _%s_ - _%s_", c[1], c[2], c[4], endtime[1])
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

	return strings.Join(rows, "\n")
}

func getPostableTwitterData() string {
	if len(updatedLocations.Locations) == 0 {
		return ""
	}

	rows := make([]string, 0)
	for _, location := range updatedLocations.Locations {
		if location.ChangeType == "REMOVED" {
			rows = append(rows, fmt.Sprintf("REMOVED: %s", location.TwitterData))
		} else {
			rows = append(rows, location.TwitterData)
		}
	}

	return strings.Join(rows, "\n\n")
}
