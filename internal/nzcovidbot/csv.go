package nzcovidbot

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// Slice of updated located
type UpdatedLocations struct {
	Locations []UpdatedRow
}

// Updated data
type UpdatedRow struct {
	ChangeDate  time.Time // To order b
	ChangeType  string    // ADDED, REMOVED, MODIFIED
	DiscordData string    // Formatted Row data
	TwitterData string    // Formatted Row data
	SlackData   string    // Formatted Row data
}

// Struct of updated locations
var updatedLocations UpdatedLocations

// cache of [exposureID]row of row data
var rowCache map[string]string

// parseCsvRow Build into struct for output later
func parseCsvRow(changeType string, data string) {
	parsedTime := parseTimeFromRow(data)

	c := parseRawRowData(data)
	if rowHasChanged(c[4], data) {
		newRow := UpdatedRow{
			ChangeDate:  parsedTime,
			ChangeType:  changeType,
			DiscordData: formatCsvDiscordRow(c),
			TwitterData: formatCsvTwitterRow(c),
			SlackData:   formatCsvSlackRow(c),
		}

		// Update row cache
		rowCache[c[4]] = data

		// Append row data
		updatedLocations.Locations = append(updatedLocations.Locations, newRow)
	}
}

// rowHasChanged - Determine if row has actually changed
func rowHasChanged(exposureId string, row string) bool {
	val, exists := rowCache[exposureId]
	if !exists {
		return true
	}

	if val != row {
		return true
	}

	return false
}

// loadRepoIntoCache - reads all CSV data and parses the rows into our cache
func loadRepoIntoCache(repoLocation string) {
	// Init our cache!
	rowCache = make(map[string]string)

	folders, err := ioutil.ReadDir(repoLocation + "/locations-of-interest")
	if err != nil {
		log.Fatal(err)
	}

	// /august-2021
	for _, f := range folders {
		if f.IsDir() {
			files, err := ioutil.ReadDir(repoLocation + "/locations-of-interest/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}

			// august-2021/locations-of-interest.csv
			for _, x := range files {
				fullLocation := repoLocation + "/locations-of-interest/" + f.Name() + "/" + x.Name()
				if strings.HasSuffix(fullLocation, ".csv") {
					loadRowsIntoCache(fullLocation)
				}
			}
		}
	}

	log.Printf("Successfully populated cache with %d entries", len(rowCache))
}

func loadRowsIntoCache(filePath string) {
	// Open the file
	csvfile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)

	// Iterate through the records
	i := 0
	for {
		// Skip header row
		if i == 0 {
			i++
			continue
		}

		// Read each record from csv
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Add to cache var
		rowCache[row[0]] = strings.Join(row, ",")
	}
}

func orderRowDataByDate() {
	sort.Slice(updatedLocations.Locations, func(i, j int) bool {
		return updatedLocations.Locations[i].ChangeDate.Before(updatedLocations.Locations[j].ChangeDate)
	})
}

// formatCsvDiscordRow Format the string to a tidy string for the interwebs
func formatCsvDiscordRow(c []string) string {
	return fmt.Sprintf("**%s** %s on _%s_ - _%s_", c[2], c[3], c[0], c[1])
}

// formatCsvTwitterRow Format the string to a tidy string for the interwebs
func formatCsvTwitterRow(c []string) string {
	return fmt.Sprintf("New Location: %s\n%s\n%s - %s\n#NZCovidTracker #NZCovid", c[2], c[3], c[0], c[1])
}

// formatCsvSlackRow Format the string to a tidy string for the interwebs
func formatCsvSlackRow(c []string) string {
	return fmt.Sprintf("*%s* %s on _%s_ - _%s_", c[2], c[3], c[0], c[1])
}

func parseTimeFromRow(data string) time.Time {
	r := csv.NewReader(strings.NewReader(data))
	r.Comma = ','
	fields, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return time.Now()
	}

	c := make([]string, 0)
	c = append(c, fields...)

	starttime := c[4]
	st, err := time.Parse("02/01/2006, 3:04 pm", starttime)
	if err != nil {
		log.Print(err)
		return time.Now()
	}

	return st
}

// Returns []string of parsed data.. starttime, endtime, name, address, ID
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

	return append(output, starttime, endtime, c[1], c[2], c[0])
}

func getPostableDiscordData() []string {
	groups := make([]string, 0)
	if len(updatedLocations.Locations) == 0 {
		return groups
	}

	rows := make([]string, 0)
	for _, location := range updatedLocations.Locations {
		if location.ChangeType == "REMOVED" {
			rows = append(rows, fmt.Sprintf("REMOVED: %s", location.DiscordData))
		} else {
			rows = append(rows, location.DiscordData)
		}

		if len(rows) > 20 {
			groups = append(groups, strings.Join(rows, "\n"))
			rows = make([]string, 0)
		}
	}

	return append(groups, strings.Join(rows, "\n"))
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
