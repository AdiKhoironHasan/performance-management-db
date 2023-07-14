package storage

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	dateString := "24-Okt-2022"
	layout := "2-Jan-2006"

	// Define a mapping of month abbreviations to their corresponding time.Month values
	monthAbbreviations := map[string]time.Month{
		"Jan": time.January,
		"Feb": time.February,
		"Mar": time.March,
		"Apr": time.April,
		"Mei": time.May,
		"Jun": time.June,
		"Jul": time.July,
		"Agt": time.August,
		"Sep": time.September,
		"Okt": time.October,
		"Nov": time.November,
		"Des": time.December,
	}

	// Parse the string into a time.Time value
	times, err := parseCustomLayout(dateString, layout, monthAbbreviations)
	if err != nil {
		log.Fatal(err)
	}

	// Output the parsed time value
	log.Println(times)

}

// Custom function to parse the string with a custom month abbreviation mapping
func parseCustomLayout(dateString, layout string, monthAbbreviations map[string]time.Month) (time.Time, error) {
	// Replace the month abbreviations with the numeric representation
	for abbreviation, month := range monthAbbreviations {
		dateString = strings.Replace(dateString, abbreviation, fmt.Sprintf("%02d", month), 1)
	}

	return time.Parse(layout, dateString)
}
