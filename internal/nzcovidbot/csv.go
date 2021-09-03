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
	FromDate        time.Time `json:"FromDate"`        // Start date
	EndDate         time.Time `json:"EndDate"`         // End date
	LocationName    string    `json:"LocationName"`    // Location Name
	LocationAddress string    `json:"LocationAddress"` // Location Address

	DiscordData string `json:"-"` // Formatted Row data
	TwitterData string `json:"-"` // Formatted Row data
	SlackData   string `json:"-"` // Formatted Row data
}

// Struct of updated locations
var updatedLocations UpdatedLocations

// cache of [exposureID]row of row data
var rowCache map[string]UpdatedRow

// parseCsvRow Build into struct for output later
func parseCsvRow(data string) {
	c, st, et := parseRawRowData(data)

	if rowHasChanged(c[4], st, et, c[2], c[3]) {
		newRow := UpdatedRow{
			FromDate:        st,
			EndDate:         et,
			LocationName:    c[2],
			LocationAddress: c[3],
			DiscordData:     formatCsvDiscordRow(c),
			TwitterData:     formatCsvTwitterRow(c),
			SlackData:       formatCsvSlackRow(c),
		}

		// Update row cache! [exposureId]UpdatedRow
		rowCache[c[4]] = newRow

		// Append row data
		updatedLocations.Locations = append(updatedLocations.Locations, newRow)
	}
}

// rowHasChanged - Determine if row has actually changed based on raw data
func rowHasChanged(exposureId string, startTime time.Time, endTime time.Time, locationName string, locationAddress string) bool {
	val, exists := rowCache[exposureId]
	if !exists {
		log.Printf("exposureId %s is new. Adding to cache", exposureId)
		return true
	}

	if val.FromDate.Unix() != startTime.Unix() {
		log.Printf("StartDate Change for %s from %s to %s", exposureId, val.FromDate.String(), startTime.String())
		return true
	}

	if val.EndDate.Unix() != endTime.Unix() {
		log.Printf("EndDate Change for %s from %s to %s", exposureId, val.EndDate.String(), endTime.String())
		return true
	}

	if !strings.EqualFold(val.LocationName, locationName) {
		log.Printf("LocationName Change for %s from %s to %s", exposureId, val.LocationName, locationName)
		return true
	}

	if !strings.EqualFold(val.LocationAddress, locationAddress) {
		log.Printf("LocationAddress Change for %s from %s to %s", exposureId, val.LocationAddress, locationAddress)
		return true
	}

	return false
}

// loadRepoIntoCache - reads all CSV data and parses the rows into our cache
func loadRepoIntoCache(repoLocation string) {
	// Init our cache!
	rowCache = make(map[string]UpdatedRow)

	// Load cache file. ELSE load files.

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
		// Read each record from csv
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Skip header row
		if i == 0 {
			i++
			continue
		}

		// Parse into our required format
		c := make([]string, 0)
		c = append(c, row...)

		st, et := parseRowTimes(c[4], c[5])

		// Build object
		newRow := UpdatedRow{
			FromDate:        st,
			EndDate:         et,
			LocationName:    c[1],
			LocationAddress: c[2],
		}

		// Add to cache
		rowCache[row[0]] = newRow
	}
}

func orderRowDataByDate() {
	sort.Slice(updatedLocations.Locations, func(i, j int) bool {
		return updatedLocations.Locations[i].FromDate.Before(updatedLocations.Locations[j].FromDate)
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

// Returns []string of parsed data.. starttime, endtime, name, address, ID
func parseRawRowData(data string) ([]string, time.Time, time.Time) {
	output := make([]string, 0)

	r := csv.NewReader(strings.NewReader(data))
	r.Comma = ','
	fields, err := r.Read()
	if err != nil {
		fmt.Println(err)
		return output, time.Now(), time.Now()
	}

	c := make([]string, 0)
	c = append(c, fields...)

	st, et := parseRowTimes(c[4], c[5])

	starttime := st.Format("Monday 2 Jan, 3:04PM")
	endtime := et.Format("3:04PM")

	return append(output, starttime, endtime, c[1], c[2], c[0]), st, et
}

func parseRowTimes(startString string, endString string) (time.Time, time.Time) {
	st, err := time.Parse("2/01/2006, 3:04 pm", startString)
	if err != nil {
		log.Print(err)
		st, err = time.Parse("2006-01-02 15:04:05", startString)
		if err != nil {
			log.Print(err)
			st = time.Now()
		}
	}

	et, err := time.Parse("2/01/2006, 3:04 pm", endString)
	if err != nil {
		log.Print(err)
		et, err = time.Parse("2006-01-02 15:04:05", endString)
		if err != nil {
			log.Print(err)
			et = time.Now()
		}
	}
	return st, et
}

func getPostableDiscordData() []string {
	groups := make([]string, 0)
	if len(updatedLocations.Locations) == 0 {
		return groups
	}

	rows := make([]string, 0)
	for _, location := range updatedLocations.Locations {
		rows = append(rows, location.DiscordData)

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
		rows = append(rows, location.SlackData)
	}

	return rows
}
