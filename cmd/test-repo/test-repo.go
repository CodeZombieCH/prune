// Simulates creating a backup every Sunday,
// starting from the first Sunday of the year 2000, ending after the first Sunday of the year 2001
//
// A backup is currently represented by a time stamped directory which contains the backups of that day

package main

import (
	"fmt"
	"os"
	"path"
	"time"
)

var (
	testBaseDirectory string
)

func main() {
	// Parse args
	args := os.Args[1:]
	if expectedArgs := 1; len(args) != expectedArgs {
		fmt.Printf("Invalid number of arguments: expected %v, got %v\n", expectedArgs, len(args))
		os.Exit(2)
	}
	testBaseDirectory = args[0]

	// Run
	if err := run(); err != nil {
		fmt.Printf("Shit hit the fan: %s", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(testBaseDirectory, os.ModePerm); err != nil {
		return err
	}

	layout := "2006-01-02T15-04-05.000Z"
	date := time.Date(2000, 01, 02, 0, 0, 0, 0, time.UTC)

	for date.Before(time.Date(2001, 01, 8, 0, 0, 0, 0, time.UTC)) {
		directoryName := date.Format(layout)
		fmt.Println(directoryName)

		if err := os.Mkdir(path.Join(testBaseDirectory, directoryName), os.ModePerm); err != nil {
			return err
		}

		date = date.AddDate(0, 0, 7)
	}

	return nil
}
