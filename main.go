package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	flag "github.com/spf13/pflag"
)

var (
	errorLogger *log.Logger

	baseDirectory string
	keepDaily     int
	keepMonthly   int
	keepYearly    int
)

func init() {
	errorLogger = log.New(os.Stderr, "", 0)

	flag.IntVarP(&keepDaily, "keep-daily", "d", -1, "number of daily files/directories to keep")
	flag.IntVarP(&keepMonthly, "keep-monthly", "m", -1, "number of monthly files/directories to keep")
	flag.IntVarP(&keepYearly, "keep-yearly", "y", -1, "number of yearly files/directories to keep")
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
	fmt.Println("keep-daily: ", keepDaily)
	fmt.Println("keep-monthly: ", keepMonthly)
	fmt.Println("keep-yearly: ", keepYearly)

	// Create config
	config := Configuration{Path: baseDirectory, KeepDaily: keepDaily, KeepMonthly: keepMonthly, KeepYearly: keepYearly}

	traverser := FileSystemTraverser{}
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
		errorLogger.Printf("%s: %t", object.Directory.Path, object.Keep)
	}
}
