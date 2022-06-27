package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

func TestParseDateRangesISODateOnly(t *testing.T) {
	// Arrange
	var args = []string{
		"2000-02-28",
		"2001-12-01...2002-01-31",
	}

	// Act
	dateRanges, err := parseDateRanges(args)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if expected, actual := 2, len(dateRanges); actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	var dateRange dateRange

	dateRange = dateRanges[0]
	if expected, actual := time.Date(2000, 2, 28, 0, 0, 0, 0, time.UTC), dateRange.From; actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
	if expected, actual := time.Date(2000, 2, 28, 0, 0, 0, 0, time.UTC), dateRange.To; actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	dateRange = dateRanges[1]
	if expected, actual := time.Date(2001, 12, 1, 0, 0, 0, 0, time.UTC), dateRange.From; actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
	if expected, actual := time.Date(2002, 1, 31, 0, 0, 0, 0, time.UTC), dateRange.To; actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestParseDateRangesISODateAndTime(t *testing.T) {
	// Arrange
	var args = []string{
		"2000-02-28T01:02:03Z",
		"2001-12-01T01:02:03Z...2002-01-31T01:02:03Z",
	}

	// Act
	dateRanges, err := parseDateRanges(args)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if expected, actual := 2, len(dateRanges); actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	var dateRange dateRange

	dateRange = dateRanges[0]
	if expected, actual := time.Date(2000, 2, 28, 1, 2, 3, 0, time.UTC), dateRange.From; actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
	if expected, actual := time.Date(2000, 2, 28, 1, 2, 3, 0, time.UTC), dateRange.To; actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	dateRange = dateRanges[1]
	if expected, actual := time.Date(2001, 12, 1, 1, 2, 3, 0, time.UTC), dateRange.From; actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
	if expected, actual := time.Date(2002, 1, 31, 1, 2, 3, 0, time.UTC), dateRange.To; actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestParseDateRangesInvalid(t *testing.T) {
	var testCases = [][]string{
		{"Monday, 02-Jan-06 15:04:05 MST"},
		{"2001-12-01..."},
		{"2001-12-01...2002-01-31...2003-01-31"},
		{"2001-13-01..."},
		{"2000-01-01...2000-13-01"},
	}

	for _, args := range testCases {
		t.Run(fmt.Sprintf("args: %v", args), func(t *testing.T) {
			_, err := parseDateRanges(args)
			var e *ArgParseError
			if err == nil || !errors.As(err, &e) {
				t.Fatal("expected error of type ArgParseError")
			}
		})
	}
}

func TestRepoISODateOnly(t *testing.T) {
	repoPath := t.TempDir()
	fmt.Println(repoPath)

	// Arrange
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"cmd",
		"--pattern",
		"%Y-%m-%d",
		repoPath,
		"2000-01-01...2000-01-31",
		"2000-02-28",
		"2000-03-01...2000-03-31",
	}

	// Act
	pattern = "%Y-%m-%d" // <========= Evil hack: override parsed and cached flags
	main()

	cmd := exec.Command("ls", "-la", repoPath)
	out, _ := cmd.CombinedOutput()
	fmt.Println(string(out))

	// Assert
	dirSpotChecks := map[string]bool{
		"2000-01-01": true,
		"2000-01-02": true,
		"2000-01-30": true,
		"2000-01-31": true,
		"2000-02-01": false,
		"2000-02-27": false,
		"2000-02-28": true,
		"2000-02-29": false,
		"2000-03-01": true,
		"2000-03-02": true,
		"2000-03-30": true,
		"2000-03-31": true,
	}
	assertExistance(repoPath, dirSpotChecks, t)
}

func TestRepoISODateAndTime(t *testing.T) {
	repoPath := t.TempDir()
	fmt.Println(repoPath)

	// Arrange
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{
		"cmd",
		repoPath,
		"2000-01-01...2000-01-31",
		"2000-02-28",
		"2000-03-01...2000-03-31",
	}

	// Act
	pattern = "%Y-%m-%dT%H-%M-%SZ" // <========= Evil hack: override parsed and cached flags
	main()

	cmd := exec.Command("ls", "-la", repoPath)
	out, _ := cmd.CombinedOutput()
	fmt.Println(string(out))

	// Assert
	dirSpotChecks := map[string]bool{
		"2000-01-01T00-00-00Z": true,
		"2000-01-02T00-00-00Z": true,
		"2000-01-30T00-00-00Z": true,
		"2000-01-31T00-00-00Z": true,
		"2000-02-01T00-00-00Z": false,
		"2000-02-27T00-00-00Z": false,
		"2000-02-28T00-00-00Z": true,
		"2000-02-29T00-00-00Z": false,
		"2000-03-01T00-00-00Z": true,
		"2000-03-02T00-00-00Z": true,
		"2000-03-30T00-00-00Z": true,
		"2000-03-31T00-00-00Z": true,
	}
	assertExistance(repoPath, dirSpotChecks, t)
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func assertExistance(directoryPath string, expectations map[string]bool, t *testing.T) {
	for dirName, expected := range expectations {
		exists := dirExists(path.Join(directoryPath, dirName))
		if expected != exists {
			if expected {
				t.Errorf("Expected directory %v to exists, but does not", dirName)
			} else {
				t.Errorf("Expected directory %v to not exists, but does", dirName)
			}
		}
	}
}
