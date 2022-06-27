// Utility to create a backup repository allowing to test pruning against
package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/itchyny/timefmt-go"
	flag "github.com/spf13/pflag"
)

const DateRangeSeperator = "..."
const PatternAlmostISO8601WithZuluTime = "%Y-%m-%dT%H-%M-%SZ"

// CLI arguments date patterns (used to parse input)
const (
	ArgumentDatePatternISO8601DateOnly    = "2006-01-02"               // "%Y-%m-%d"
	ArgumentDatePatternISO8601DateAndTime = "2006-01-02T15:04:05Z0700" // "%Y-%m-%dT%H:%M:%S%z"
)

type ArgParseError struct {
	Position   int
	Value      string
	InnerError error
}

func (e *ArgParseError) Error() string {
	return fmt.Sprintf("Invalid argument #%d '%v': %v", e.Position, e.Value, e.InnerError)
}

func (e *ArgParseError) Unwrap() error { return e.InnerError }

type configuration struct {
	BaseDirectory string
	DateRanges    []dateRange
	Pattern       string
}

func run(config configuration) error {
	if err := os.RemoveAll(config.BaseDirectory); err != nil {
		return err
	}

	if err := os.MkdirAll(config.BaseDirectory, os.ModePerm); err != nil {
		return err
	}

	for _, v := range config.DateRanges {
		for date := v.From; !date.After(v.To); date = date.AddDate(0, 0, 1) {
			directoryName := timefmt.Format(date, config.Pattern)
			directoryPath := path.Join(config.BaseDirectory, directoryName)

			// Create directory
			if err := os.Mkdir(directoryPath, os.ModePerm); err != nil {
				return err
			}

			// Create test files inside directory
			files := []string{
				"backup-" + timefmt.Format(date, PatternAlmostISO8601WithZuluTime) + ".tar.gz",
				"backup-" + timefmt.Format(date, PatternAlmostISO8601WithZuluTime) + ".tar.gz.sha256sum",
			}
			for _, v := range files {
				err := os.WriteFile(path.Join(directoryPath, v), []byte(v), 0644)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func parseDateString(s string) (time.Time, error) {
	var date time.Time
	var err error

	// Try to parse as date only format
	date, err = time.Parse(ArgumentDatePatternISO8601DateOnly, s)

	// If error, try to parse as date time format
	if err != nil {
		date, err = time.Parse(ArgumentDatePatternISO8601DateAndTime, s)
	}

	// If error, give up
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func parseDateRanges(args []string) ([]dateRange, error) {
	dateRanges := []dateRange{}
	for i, v := range args {
		if strings.Contains(v, DateRangeSeperator) {
			dates := strings.Split(v, DateRangeSeperator)
			if len(dates) != 2 {
				return nil, &ArgParseError{Position: i, Value: v}
			}
			from, err := parseDateString(dates[0])
			if err != nil {
				return nil, &ArgParseError{Position: i, Value: v, InnerError: err}
			}
			to, err := parseDateString(dates[1])
			if err != nil {
				return nil, &ArgParseError{Position: i, Value: v, InnerError: err}
			}

			dateRanges = append(dateRanges, dateRange{From: from, To: to})
		} else {
			date, err := parseDateString(v)
			if err != nil {
				return nil, &ArgParseError{Position: i, Value: v, InnerError: err}
			}
			dateRanges = append(dateRanges, dateRange{From: date, To: date})
		}
	}
	return dateRanges, nil
}

type dateRange struct {
	From time.Time
	To   time.Time
}

var pattern string

func init() {
	flag.StringVarP(&pattern, "pattern", "p", PatternAlmostISO8601WithZuluTime, "strptime pattern used to parse the date from the name of the timestamped directory")
}

func main() {
	config := configuration{}

	// Parse flags
	flag.Parse()
	config.Pattern = pattern

	// Parse args
	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("Invalid number of arguments: provide at least two arguments")
		os.Exit(2)
	}

	// Argument #1: directory
	config.BaseDirectory = args[0]

	// Arguments #2-#n: date (date ranges)
	dateRanges, err := parseDateRanges(args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	config.DateRanges = dateRanges

	// Run
	if err := run(config); err != nil {
		fmt.Printf("Shit hit the fan: %s", err)
		os.Exit(1)
	}
}
