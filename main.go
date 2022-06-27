// Copyright 2022 Marc-André Bühler
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"os"
	"sort"

	flag "github.com/spf13/pflag"
)

var (
	logger      *log.Logger
	errorLogger *log.Logger

	baseDirectory string
	verbose       bool
	keepDaily     int
	keepMonthly   int
	keepYearly    int
	pattern       string
)

func init() {
	logger = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)

	flag.BoolVarP(&verbose, "verbose", "v", false, "verbose flag")

	flag.IntVarP(&keepDaily, "keep-daily", "d", -1, "number of daily files/directories to keep")
	flag.IntVarP(&keepMonthly, "keep-monthly", "m", -1, "number of monthly files/directories to keep")
	flag.IntVarP(&keepYearly, "keep-yearly", "y", -1, "number of yearly files/directories to keep")

	// TODO: evaluate sane default (if a default makes sense at all)
	flag.StringVarP(&pattern, "pattern", "p", PatternAlmostISO8601DateAndTime, "strptime pattern used to parse the date from the name of the timestamped directory")
}

func main() {
	// Parse
	flag.Parse()

	if flag.NArg() != 1 {
		errorLogger.Printf("To many arguments")
		os.Exit(2) // Aligns with pflag "ExitOnError will call os.Exit(2) if an error is found when parsing"
	}
	baseDirectory = flag.Args()[0]

	// Validate

	// Run
	if err := run(); err != nil {
		errorLogger.Printf("Shit hit the fan: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if verbose {
		logger.Printf("keep-daily: %v, keep-monthly: %v, keep-yearly: %v", keepDaily, keepMonthly, keepYearly)
	}

	// TODO: validate pattern

	// Create config
	config := Configuration{
		Path:        baseDirectory,
		Pattern:     pattern,
		KeepDaily:   keepDaily,
		KeepMonthly: keepMonthly,
		KeepYearly:  keepYearly,
	}

	traverser := FileSystemTraverser{Pattern: pattern}
	objects, err := traverser.GetObjects(baseDirectory)
	if err != nil {
		errorLogger.Printf("Failed to retrieve directories")
		return err
	}

	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(objects)
	if err != nil {
		errorLogger.Printf("Failed to calculate directories to prune")
		return err
	}

	printSorted(pruneResult.Objects)

	if verbose {
		printStats(pruneResult)
	}

	return err
}

func printSorted(objects map[string]*PruneCandidate) {
	keys := make([]string, 0, len(objects))
	for k := range objects {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		object := objects[k]

		if verbose {
			var operation string
			switch object.Keep {
			case true:
				operation = "keep"
			case false:
				operation = "prune"
			}
			errorLogger.Printf("%s: %s\n", object.Directory.Path, operation)
		} else {
			// Print only directories to prune
			if !object.Keep {
				logger.Println(object.Directory.Path)
			}
		}
	}
}

func printStats(result PruneResult) {
	logger.Printf("Total count: keep: %v, prune: %v\n", len(result.ToKeep), len(result.ToPrune))
}
