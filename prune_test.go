package main

import (
	"io/fs"
	"path"
	"testing"
	"time"
)

var testBaseDirectory string = "/foo/bar"

func TestPruneEmpty(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 7}
	testDirectories := []TestObject{}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 0; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}
	if expected := 0; len(pruneResult.ToPrune) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToPrune), expected)
	}
}

func TestPruneNothing(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: NoPrune, KeepMonthly: NoPrune, KeepYearly: NoPrune}
	testDirectories := []TestObject{
		{"2000-01-01T00-00-00.000Z", true},
		{"2000-01-02T00-00-00.000Z", true},
		{"2000-01-03T00-00-00.000Z", true},
		{"2000-01-04T00-00-00.000Z", true},
		{"2000-01-05T00-00-00.000Z", true},
		{"2000-01-06T00-00-00.000Z", true},
		{"2000-01-07T00-00-00.000Z", true},
		{"2000-01-08T00-00-00.000Z", true},
		{"2000-01-09T00-00-00.000Z", true},
		{"2000-01-10T00-00-00.000Z", true},
		{"2000-01-11T00-00-00.000Z", true},
		{"2000-01-12T00-00-00.000Z", true},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 12; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}
	if expected := 0; len(pruneResult.ToPrune) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToPrune), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneEverything(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 0, KeepMonthly: 0, KeepYearly: 0}
	testDirectories := []TestObject{
		{"2000-01-01T00-00-00.000Z", false},
		{"2000-01-02T00-00-00.000Z", false},
		{"2000-01-03T00-00-00.000Z", false},
		{"2000-01-04T00-00-00.000Z", false},
		{"2000-01-05T00-00-00.000Z", false},
		{"2000-01-06T00-00-00.000Z", false},
		{"2000-01-07T00-00-00.000Z", false},
		{"2000-01-08T00-00-00.000Z", false},
		{"2000-01-09T00-00-00.000Z", false},
		{"2000-01-10T00-00-00.000Z", false},
		{"2000-01-11T00-00-00.000Z", false},
		{"2000-01-12T00-00-00.000Z", false},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 0; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}
	if expected := 12; len(pruneResult.ToPrune) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToPrune), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneEverythingUsingKeepDaily(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 0}
	testDirectories := []TestObject{
		{"2000-01-01T00-00-00.000Z", false},
		{"2000-01-02T00-00-00.000Z", false},
		{"2000-01-03T00-00-00.000Z", false},
		{"2000-01-04T00-00-00.000Z", false},
		{"2000-01-05T00-00-00.000Z", false},
		{"2000-01-06T00-00-00.000Z", false},
		{"2000-01-07T00-00-00.000Z", false},
		{"2000-01-08T00-00-00.000Z", false},
		{"2000-01-09T00-00-00.000Z", false},
		{"2000-01-10T00-00-00.000Z", false},
		{"2000-01-11T00-00-00.000Z", false},
		{"2000-01-12T00-00-00.000Z", false},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 0; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}
	if expected := 12; len(pruneResult.ToPrune) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToPrune), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneWithMultiplePerDay(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 4}
	testDirectories := []TestObject{
		{"2000-01-01T01-00-00.000Z", false},
		{"2000-01-01T02-00-00.000Z", true},

		{"2000-01-02T01-01-00.000Z", false},
		{"2000-01-02T01-02-00.000Z", false},
		{"2000-01-02T01-03-00.000Z", true},

		{"2000-01-03T01-01-01.000Z", false},
		{"2000-01-03T01-01-02.000Z", false},
		{"2000-01-03T01-01-03.000Z", false},
		{"2000-01-03T01-01-04.000Z", true},

		{"2000-01-04T01-01-01.001Z", false},
		{"2000-01-04T01-01-01.002Z", false},
		{"2000-01-04T01-01-01.003Z", false},
		{"2000-01-04T01-01-01.004Z", false},
		{"2000-01-04T01-01-01.005Z", true},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 4; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}
	if expected := 10; len(pruneResult.ToPrune) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToPrune), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneDaily(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 7, KeepMonthly: NoPrune, KeepYearly: NoPrune}
	testDirectories := []TestObject{
		{"2000-01-01T00-00-00.000Z", false},
		{"2000-01-02T00-00-00.000Z", false},
		{"2000-01-03T00-00-00.000Z", false},
		{"2000-01-04T00-00-00.000Z", false},
		{"2000-01-05T00-00-00.000Z", false},
		{"2000-01-06T00-00-00.000Z", true},
		{"2000-01-07T00-00-00.000Z", true},
		{"2000-01-08T00-00-00.000Z", true},
		{"2000-01-09T00-00-00.000Z", true},
		{"2000-01-10T00-00-00.000Z", true},
		{"2000-01-11T00-00-00.000Z", true},
		{"2000-01-12T00-00-00.000Z", true},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 7; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}
	if expected := 5; len(pruneResult.ToPrune) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToPrune), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneMonthly(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: NoPrune, KeepMonthly: 12, KeepYearly: NoPrune}
	testDirectories := []TestObject{
		{"2000-01-01T00-00-00.000Z", false},
		{"2000-02-01T00-00-00.000Z", false},
		{"2000-03-01T00-00-00.000Z", false},
		{"2000-04-01T00-00-00.000Z", false},
		{"2000-05-01T00-00-00.000Z", false},
		{"2000-06-01T00-00-00.000Z", false},
		{"2000-07-01T00-00-00.000Z", true},
		{"2000-08-01T00-00-00.000Z", true},
		{"2000-09-01T00-00-00.000Z", true},
		{"2000-10-01T00-00-00.000Z", true},
		{"2000-11-01T00-00-00.000Z", true},
		{"2000-12-01T00-00-00.000Z", true},
		{"2001-01-01T00-00-00.000Z", true},
		{"2001-02-01T00-00-00.000Z", true},
		{"2001-03-01T00-00-00.000Z", true},
		{"2001-04-01T00-00-00.000Z", true},
		{"2001-05-01T00-00-00.000Z", true},
		{"2001-06-01T00-00-00.000Z", true},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 12; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}
	if expected := 6; len(pruneResult.ToPrune) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToPrune), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneYearly(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: NoPrune, KeepMonthly: NoPrune, KeepYearly: 10}
	testDirectories := []TestObject{
		{"2000-01-01T00-00-00.000Z", false},
		{"2001-01-01T00-00-00.000Z", false},
		{"2002-01-01T00-00-00.000Z", false},
		{"2003-01-01T00-00-00.000Z", false},
		{"2004-01-01T00-00-00.000Z", false},
		{"2005-01-01T00-00-00.000Z", true},
		{"2006-01-01T00-00-00.000Z", true},
		{"2007-01-01T00-00-00.000Z", true},
		{"2008-01-01T00-00-00.000Z", true},
		{"2009-01-01T00-00-00.000Z", true},
		{"2010-01-01T00-00-00.000Z", true},
		{"2011-01-01T00-00-00.000Z", true},
		{"2012-01-01T00-00-00.000Z", true},
		{"2013-01-01T00-00-00.000Z", true},
		{"2014-01-01T00-00-00.000Z", true},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 10; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}
	if expected := 5; len(pruneResult.ToPrune) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToPrune), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneDailyAndMonthly(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 7, KeepMonthly: 3}
	testDirectories := []TestObject{
		{"2000-01-03T00-00-00.000Z", false},
		{"2000-01-10T00-00-00.000Z", false},
		{"2000-01-17T00-00-00.000Z", false},
		{"2000-01-24T00-00-00.000Z", false},
		{"2000-01-31T00-00-00.000Z", true},

		{"2000-02-07T00-00-00.000Z", false},
		{"2000-02-14T00-00-00.000Z", false},
		{"2000-02-21T00-00-00.000Z", false},
		{"2000-02-28T00-00-00.000Z", true},

		{"2000-03-06T00-00-00.000Z", false},
		{"2000-03-13T00-00-00.000Z", false},
		{"2000-03-20T00-00-00.000Z", false},
		{"2000-03-27T00-00-00.000Z", true},

		{"2000-04-03T00-00-00.000Z", true},
		{"2000-04-10T00-00-00.000Z", true},
		{"2000-04-17T00-00-00.000Z", true},
		{"2000-04-24T00-00-00.000Z", true},

		{"2000-05-01T00-00-00.000Z", true},
		{"2000-05-08T00-00-00.000Z", true},
		{"2000-05-15T00-00-00.000Z", true},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 10; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneOldest(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 3, KeepMonthly: 2, KeepYearly: 1}
	testDirectories := []TestObject{
		{"2000-01-01T00-00-00.000Z", true},
		{"2000-04-01T00-00-00.000Z", false},
		{"2000-07-01T00-00-00.000Z", false},
		{"2000-10-01T00-00-00.000Z", true},

		{"2001-01-01T00-00-00.000Z", true},
		{"2001-04-01T00-00-00.000Z", true},
		{"2001-07-01T00-00-00.000Z", true},
		{"2001-10-01T00-00-00.000Z", true},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	if expected := 6; len(pruneResult.ToKeep) != expected {
		t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	}

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

func TestPruneTotalCountWithAddedBackups(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 3, KeepMonthly: 2, KeepYearly: 1}
	layout := "2006-01-02T15-04-05.000Z"

	// Act
	startDate := time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 400; i++ {
		// Prepare test objects
		var currentDate time.Time
		testDirectories := []TestObject{}
		for j := 0; j <= i; j++ {
			currentDate = startDate.AddDate(0, 0, j)
			testDirectories = append(testDirectories, TestObject{currentDate.Format(layout), true /*not relevant*/})
		}
		entries := createEntries(testDirectories, t)

		t.Logf("to %v: %v", startDate.AddDate(0, 0, i), len(entries))
		// t.Logf("entries: %v", entries)

		prune := NewPrune(config)
		pruneResult, err := prune.Calculate(entries)
		if err != nil {
			t.Fatalf("Failed to calculate directories to prune: %s", err)
		}

		// Assert
		var expected int
		switch {
		case currentDate.Before(time.Date(2000, 01, 03, 0, 0, 0, 0, time.UTC)) ||
			currentDate == time.Date(2000, 01, 03, 0, 0, 0, 0, time.UTC):
			expected = len(entries)
		case currentDate.Before(time.Date(2000, 02, 03, 0, 0, 0, 0, time.UTC)):
			expected = 4
		case currentDate.Before(time.Date(2000, 03, 03, 0, 0, 0, 0, time.UTC)):
			expected = 5
		default:
			expected = 6
		}

		if len(pruneResult.ToKeep) != expected {
			t.Errorf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
		}
	}
}

// Assume it is 2016-01-01, today's backup has not yet been made, you have
// created at least one backup on each day in 2015 except on 2015-12-19 (no
// backup made on that day), and you started backing up with borg on
// 2015-01-01.
func TestPruneBorgExample(t *testing.T) {
	// Arrange
	config := Configuration{Path: testBaseDirectory, KeepDaily: 14, KeepMonthly: 6, KeepYearly: 1}
	testDirectories := []TestObject{
		{"2015-01-01T00-00-00.000Z", true},
		{"2015-01-02T00-00-00.000Z", false},
		{"2015-01-03T00-00-00.000Z", false},
		{"2015-01-04T00-00-00.000Z", false},
		{"2015-01-05T00-00-00.000Z", false},
		{"2015-01-06T00-00-00.000Z", false},
		{"2015-01-07T00-00-00.000Z", false},
		{"2015-01-08T00-00-00.000Z", false},
		{"2015-01-09T00-00-00.000Z", false},
		{"2015-01-10T00-00-00.000Z", false},
		{"2015-01-11T00-00-00.000Z", false},
		{"2015-01-12T00-00-00.000Z", false},
		{"2015-01-13T00-00-00.000Z", false},
		{"2015-01-14T00-00-00.000Z", false},
		{"2015-01-15T00-00-00.000Z", false},
		{"2015-01-16T00-00-00.000Z", false},
		{"2015-01-17T00-00-00.000Z", false},
		{"2015-01-18T00-00-00.000Z", false},
		{"2015-01-19T00-00-00.000Z", false},
		{"2015-01-20T00-00-00.000Z", false},
		{"2015-01-21T00-00-00.000Z", false},
		{"2015-01-22T00-00-00.000Z", false},
		{"2015-01-23T00-00-00.000Z", false},
		{"2015-01-24T00-00-00.000Z", false},
		{"2015-01-25T00-00-00.000Z", false},
		{"2015-01-26T00-00-00.000Z", false},
		{"2015-01-27T00-00-00.000Z", false},
		{"2015-01-28T00-00-00.000Z", false},
		{"2015-01-29T00-00-00.000Z", false},
		{"2015-01-30T00-00-00.000Z", false},
		{"2015-01-31T00-00-00.000Z", false},
		{"2015-02-01T00-00-00.000Z", false},
		{"2015-02-02T00-00-00.000Z", false},
		{"2015-02-03T00-00-00.000Z", false},
		{"2015-02-04T00-00-00.000Z", false},
		{"2015-02-05T00-00-00.000Z", false},
		{"2015-02-06T00-00-00.000Z", false},
		{"2015-02-07T00-00-00.000Z", false},
		{"2015-02-08T00-00-00.000Z", false},
		{"2015-02-09T00-00-00.000Z", false},
		{"2015-02-10T00-00-00.000Z", false},
		{"2015-02-11T00-00-00.000Z", false},
		{"2015-02-12T00-00-00.000Z", false},
		{"2015-02-13T00-00-00.000Z", false},
		{"2015-02-14T00-00-00.000Z", false},
		{"2015-02-15T00-00-00.000Z", false},
		{"2015-02-16T00-00-00.000Z", false},
		{"2015-02-17T00-00-00.000Z", false},
		{"2015-02-18T00-00-00.000Z", false},
		{"2015-02-19T00-00-00.000Z", false},
		{"2015-02-20T00-00-00.000Z", false},
		{"2015-02-21T00-00-00.000Z", false},
		{"2015-02-22T00-00-00.000Z", false},
		{"2015-02-23T00-00-00.000Z", false},
		{"2015-02-24T00-00-00.000Z", false},
		{"2015-02-25T00-00-00.000Z", false},
		{"2015-02-26T00-00-00.000Z", false},
		{"2015-02-27T00-00-00.000Z", false},
		{"2015-02-28T00-00-00.000Z", false},
		{"2015-03-01T00-00-00.000Z", false},
		{"2015-03-02T00-00-00.000Z", false},
		{"2015-03-03T00-00-00.000Z", false},
		{"2015-03-04T00-00-00.000Z", false},
		{"2015-03-05T00-00-00.000Z", false},
		{"2015-03-06T00-00-00.000Z", false},
		{"2015-03-07T00-00-00.000Z", false},
		{"2015-03-08T00-00-00.000Z", false},
		{"2015-03-09T00-00-00.000Z", false},
		{"2015-03-10T00-00-00.000Z", false},
		{"2015-03-11T00-00-00.000Z", false},
		{"2015-03-12T00-00-00.000Z", false},
		{"2015-03-13T00-00-00.000Z", false},
		{"2015-03-14T00-00-00.000Z", false},
		{"2015-03-15T00-00-00.000Z", false},
		{"2015-03-16T00-00-00.000Z", false},
		{"2015-03-17T00-00-00.000Z", false},
		{"2015-03-18T00-00-00.000Z", false},
		{"2015-03-19T00-00-00.000Z", false},
		{"2015-03-20T00-00-00.000Z", false},
		{"2015-03-21T00-00-00.000Z", false},
		{"2015-03-22T00-00-00.000Z", false},
		{"2015-03-23T00-00-00.000Z", false},
		{"2015-03-24T00-00-00.000Z", false},
		{"2015-03-25T00-00-00.000Z", false},
		{"2015-03-26T00-00-00.000Z", false},
		{"2015-03-27T00-00-00.000Z", false},
		{"2015-03-28T00-00-00.000Z", false},
		{"2015-03-29T00-00-00.000Z", false},
		{"2015-03-30T00-00-00.000Z", false},
		{"2015-03-31T00-00-00.000Z", false},
		{"2015-04-01T00-00-00.000Z", false},
		{"2015-04-02T00-00-00.000Z", false},
		{"2015-04-03T00-00-00.000Z", false},
		{"2015-04-04T00-00-00.000Z", false},
		{"2015-04-05T00-00-00.000Z", false},
		{"2015-04-06T00-00-00.000Z", false},
		{"2015-04-07T00-00-00.000Z", false},
		{"2015-04-08T00-00-00.000Z", false},
		{"2015-04-09T00-00-00.000Z", false},
		{"2015-04-10T00-00-00.000Z", false},
		{"2015-04-11T00-00-00.000Z", false},
		{"2015-04-12T00-00-00.000Z", false},
		{"2015-04-13T00-00-00.000Z", false},
		{"2015-04-14T00-00-00.000Z", false},
		{"2015-04-15T00-00-00.000Z", false},
		{"2015-04-16T00-00-00.000Z", false},
		{"2015-04-17T00-00-00.000Z", false},
		{"2015-04-18T00-00-00.000Z", false},
		{"2015-04-19T00-00-00.000Z", false},
		{"2015-04-20T00-00-00.000Z", false},
		{"2015-04-21T00-00-00.000Z", false},
		{"2015-04-22T00-00-00.000Z", false},
		{"2015-04-23T00-00-00.000Z", false},
		{"2015-04-24T00-00-00.000Z", false},
		{"2015-04-25T00-00-00.000Z", false},
		{"2015-04-26T00-00-00.000Z", false},
		{"2015-04-27T00-00-00.000Z", false},
		{"2015-04-28T00-00-00.000Z", false},
		{"2015-04-29T00-00-00.000Z", false},
		{"2015-04-30T00-00-00.000Z", false},
		{"2015-05-01T00-00-00.000Z", false},
		{"2015-05-02T00-00-00.000Z", false},
		{"2015-05-03T00-00-00.000Z", false},
		{"2015-05-04T00-00-00.000Z", false},
		{"2015-05-05T00-00-00.000Z", false},
		{"2015-05-06T00-00-00.000Z", false},
		{"2015-05-07T00-00-00.000Z", false},
		{"2015-05-08T00-00-00.000Z", false},
		{"2015-05-09T00-00-00.000Z", false},
		{"2015-05-10T00-00-00.000Z", false},
		{"2015-05-11T00-00-00.000Z", false},
		{"2015-05-12T00-00-00.000Z", false},
		{"2015-05-13T00-00-00.000Z", false},
		{"2015-05-14T00-00-00.000Z", false},
		{"2015-05-15T00-00-00.000Z", false},
		{"2015-05-16T00-00-00.000Z", false},
		{"2015-05-17T00-00-00.000Z", false},
		{"2015-05-18T00-00-00.000Z", false},
		{"2015-05-19T00-00-00.000Z", false},
		{"2015-05-20T00-00-00.000Z", false},
		{"2015-05-21T00-00-00.000Z", false},
		{"2015-05-22T00-00-00.000Z", false},
		{"2015-05-23T00-00-00.000Z", false},
		{"2015-05-24T00-00-00.000Z", false},
		{"2015-05-25T00-00-00.000Z", false},
		{"2015-05-26T00-00-00.000Z", false},
		{"2015-05-27T00-00-00.000Z", false},
		{"2015-05-28T00-00-00.000Z", false},
		{"2015-05-29T00-00-00.000Z", false},
		{"2015-05-30T00-00-00.000Z", false},
		{"2015-05-31T00-00-00.000Z", false},
		{"2015-06-01T00-00-00.000Z", false},
		{"2015-06-02T00-00-00.000Z", false},
		{"2015-06-03T00-00-00.000Z", false},
		{"2015-06-04T00-00-00.000Z", false},
		{"2015-06-05T00-00-00.000Z", false},
		{"2015-06-06T00-00-00.000Z", false},
		{"2015-06-07T00-00-00.000Z", false},
		{"2015-06-08T00-00-00.000Z", false},
		{"2015-06-09T00-00-00.000Z", false},
		{"2015-06-10T00-00-00.000Z", false},
		{"2015-06-11T00-00-00.000Z", false},
		{"2015-06-12T00-00-00.000Z", false},
		{"2015-06-13T00-00-00.000Z", false},
		{"2015-06-14T00-00-00.000Z", false},
		{"2015-06-15T00-00-00.000Z", false},
		{"2015-06-16T00-00-00.000Z", false},
		{"2015-06-17T00-00-00.000Z", false},
		{"2015-06-18T00-00-00.000Z", false},
		{"2015-06-19T00-00-00.000Z", false},
		{"2015-06-20T00-00-00.000Z", false},
		{"2015-06-21T00-00-00.000Z", false},
		{"2015-06-22T00-00-00.000Z", false},
		{"2015-06-23T00-00-00.000Z", false},
		{"2015-06-24T00-00-00.000Z", false},
		{"2015-06-25T00-00-00.000Z", false},
		{"2015-06-26T00-00-00.000Z", false},
		{"2015-06-27T00-00-00.000Z", false},
		{"2015-06-28T00-00-00.000Z", false},
		{"2015-06-29T00-00-00.000Z", false},
		{"2015-06-30T00-00-00.000Z", true},
		{"2015-07-01T00-00-00.000Z", false},
		{"2015-07-02T00-00-00.000Z", false},
		{"2015-07-03T00-00-00.000Z", false},
		{"2015-07-04T00-00-00.000Z", false},
		{"2015-07-05T00-00-00.000Z", false},
		{"2015-07-06T00-00-00.000Z", false},
		{"2015-07-07T00-00-00.000Z", false},
		{"2015-07-08T00-00-00.000Z", false},
		{"2015-07-09T00-00-00.000Z", false},
		{"2015-07-10T00-00-00.000Z", false},
		{"2015-07-11T00-00-00.000Z", false},
		{"2015-07-12T00-00-00.000Z", false},
		{"2015-07-13T00-00-00.000Z", false},
		{"2015-07-14T00-00-00.000Z", false},
		{"2015-07-15T00-00-00.000Z", false},
		{"2015-07-16T00-00-00.000Z", false},
		{"2015-07-17T00-00-00.000Z", false},
		{"2015-07-18T00-00-00.000Z", false},
		{"2015-07-19T00-00-00.000Z", false},
		{"2015-07-20T00-00-00.000Z", false},
		{"2015-07-21T00-00-00.000Z", false},
		{"2015-07-22T00-00-00.000Z", false},
		{"2015-07-23T00-00-00.000Z", false},
		{"2015-07-24T00-00-00.000Z", false},
		{"2015-07-25T00-00-00.000Z", false},
		{"2015-07-26T00-00-00.000Z", false},
		{"2015-07-27T00-00-00.000Z", false},
		{"2015-07-28T00-00-00.000Z", false},
		{"2015-07-29T00-00-00.000Z", false},
		{"2015-07-30T00-00-00.000Z", false},
		{"2015-07-31T00-00-00.000Z", true},
		{"2015-08-01T00-00-00.000Z", false},
		{"2015-08-02T00-00-00.000Z", false},
		{"2015-08-03T00-00-00.000Z", false},
		{"2015-08-04T00-00-00.000Z", false},
		{"2015-08-05T00-00-00.000Z", false},
		{"2015-08-06T00-00-00.000Z", false},
		{"2015-08-07T00-00-00.000Z", false},
		{"2015-08-08T00-00-00.000Z", false},
		{"2015-08-09T00-00-00.000Z", false},
		{"2015-08-10T00-00-00.000Z", false},
		{"2015-08-11T00-00-00.000Z", false},
		{"2015-08-12T00-00-00.000Z", false},
		{"2015-08-13T00-00-00.000Z", false},
		{"2015-08-14T00-00-00.000Z", false},
		{"2015-08-15T00-00-00.000Z", false},
		{"2015-08-16T00-00-00.000Z", false},
		{"2015-08-17T00-00-00.000Z", false},
		{"2015-08-18T00-00-00.000Z", false},
		{"2015-08-19T00-00-00.000Z", false},
		{"2015-08-20T00-00-00.000Z", false},
		{"2015-08-21T00-00-00.000Z", false},
		{"2015-08-22T00-00-00.000Z", false},
		{"2015-08-23T00-00-00.000Z", false},
		{"2015-08-24T00-00-00.000Z", false},
		{"2015-08-25T00-00-00.000Z", false},
		{"2015-08-26T00-00-00.000Z", false},
		{"2015-08-27T00-00-00.000Z", false},
		{"2015-08-28T00-00-00.000Z", false},
		{"2015-08-29T00-00-00.000Z", false},
		{"2015-08-30T00-00-00.000Z", false},
		{"2015-08-31T00-00-00.000Z", true},
		{"2015-09-01T00-00-00.000Z", false},
		{"2015-09-02T00-00-00.000Z", false},
		{"2015-09-03T00-00-00.000Z", false},
		{"2015-09-04T00-00-00.000Z", false},
		{"2015-09-05T00-00-00.000Z", false},
		{"2015-09-06T00-00-00.000Z", false},
		{"2015-09-07T00-00-00.000Z", false},
		{"2015-09-08T00-00-00.000Z", false},
		{"2015-09-09T00-00-00.000Z", false},
		{"2015-09-10T00-00-00.000Z", false},
		{"2015-09-11T00-00-00.000Z", false},
		{"2015-09-12T00-00-00.000Z", false},
		{"2015-09-13T00-00-00.000Z", false},
		{"2015-09-14T00-00-00.000Z", false},
		{"2015-09-15T00-00-00.000Z", false},
		{"2015-09-16T00-00-00.000Z", false},
		{"2015-09-17T00-00-00.000Z", false},
		{"2015-09-18T00-00-00.000Z", false},
		{"2015-09-19T00-00-00.000Z", false},
		{"2015-09-20T00-00-00.000Z", false},
		{"2015-09-21T00-00-00.000Z", false},
		{"2015-09-22T00-00-00.000Z", false},
		{"2015-09-23T00-00-00.000Z", false},
		{"2015-09-24T00-00-00.000Z", false},
		{"2015-09-25T00-00-00.000Z", false},
		{"2015-09-26T00-00-00.000Z", false},
		{"2015-09-27T00-00-00.000Z", false},
		{"2015-09-28T00-00-00.000Z", false},
		{"2015-09-29T00-00-00.000Z", false},
		{"2015-09-30T00-00-00.000Z", true},
		{"2015-10-01T00-00-00.000Z", false},
		{"2015-10-02T00-00-00.000Z", false},
		{"2015-10-03T00-00-00.000Z", false},
		{"2015-10-04T00-00-00.000Z", false},
		{"2015-10-05T00-00-00.000Z", false},
		{"2015-10-06T00-00-00.000Z", false},
		{"2015-10-07T00-00-00.000Z", false},
		{"2015-10-08T00-00-00.000Z", false},
		{"2015-10-09T00-00-00.000Z", false},
		{"2015-10-10T00-00-00.000Z", false},
		{"2015-10-11T00-00-00.000Z", false},
		{"2015-10-12T00-00-00.000Z", false},
		{"2015-10-13T00-00-00.000Z", false},
		{"2015-10-14T00-00-00.000Z", false},
		{"2015-10-15T00-00-00.000Z", false},
		{"2015-10-16T00-00-00.000Z", false},
		{"2015-10-17T00-00-00.000Z", false},
		{"2015-10-18T00-00-00.000Z", false},
		{"2015-10-19T00-00-00.000Z", false},
		{"2015-10-20T00-00-00.000Z", false},
		{"2015-10-21T00-00-00.000Z", false},
		{"2015-10-22T00-00-00.000Z", false},
		{"2015-10-23T00-00-00.000Z", false},
		{"2015-10-24T00-00-00.000Z", false},
		{"2015-10-25T00-00-00.000Z", false},
		{"2015-10-26T00-00-00.000Z", false},
		{"2015-10-27T00-00-00.000Z", false},
		{"2015-10-28T00-00-00.000Z", false},
		{"2015-10-29T00-00-00.000Z", false},
		{"2015-10-30T00-00-00.000Z", false},
		{"2015-10-31T00-00-00.000Z", true},
		{"2015-11-01T00-00-00.000Z", false},
		{"2015-11-02T00-00-00.000Z", false},
		{"2015-11-03T00-00-00.000Z", false},
		{"2015-11-04T00-00-00.000Z", false},
		{"2015-11-05T00-00-00.000Z", false},
		{"2015-11-06T00-00-00.000Z", false},
		{"2015-11-07T00-00-00.000Z", false},
		{"2015-11-08T00-00-00.000Z", false},
		{"2015-11-09T00-00-00.000Z", false},
		{"2015-11-10T00-00-00.000Z", false},
		{"2015-11-11T00-00-00.000Z", false},
		{"2015-11-12T00-00-00.000Z", false},
		{"2015-11-13T00-00-00.000Z", false},
		{"2015-11-14T00-00-00.000Z", false},
		{"2015-11-15T00-00-00.000Z", false},
		{"2015-11-16T00-00-00.000Z", false},
		{"2015-11-17T00-00-00.000Z", false},
		{"2015-11-18T00-00-00.000Z", false},
		{"2015-11-19T00-00-00.000Z", false},
		{"2015-11-20T00-00-00.000Z", false},
		{"2015-11-21T00-00-00.000Z", false},
		{"2015-11-22T00-00-00.000Z", false},
		{"2015-11-23T00-00-00.000Z", false},
		{"2015-11-24T00-00-00.000Z", false},
		{"2015-11-25T00-00-00.000Z", false},
		{"2015-11-26T00-00-00.000Z", false},
		{"2015-11-27T00-00-00.000Z", false},
		{"2015-11-28T00-00-00.000Z", false},
		{"2015-11-29T00-00-00.000Z", false},
		{"2015-11-30T00-00-00.000Z", true},
		{"2015-12-01T00-00-00.000Z", false},
		{"2015-12-02T00-00-00.000Z", false},
		{"2015-12-03T00-00-00.000Z", false},
		{"2015-12-04T00-00-00.000Z", false},
		{"2015-12-05T00-00-00.000Z", false},
		{"2015-12-06T00-00-00.000Z", false},
		{"2015-12-07T00-00-00.000Z", false},
		{"2015-12-08T00-00-00.000Z", false},
		{"2015-12-09T00-00-00.000Z", false},
		{"2015-12-10T00-00-00.000Z", false},
		{"2015-12-11T00-00-00.000Z", false},
		{"2015-12-12T00-00-00.000Z", false},
		{"2015-12-13T00-00-00.000Z", false},
		{"2015-12-14T00-00-00.000Z", false},
		{"2015-12-15T00-00-00.000Z", false},
		{"2015-12-16T00-00-00.000Z", false},
		{"2015-12-17T00-00-00.000Z", true},
		{"2015-12-18T00-00-00.000Z", true},
		//{"2015-12-19T00-00-00.000Z", false},
		{"2015-12-20T00-00-00.000Z", true},
		{"2015-12-21T00-00-00.000Z", true},
		{"2015-12-22T00-00-00.000Z", true},
		{"2015-12-23T00-00-00.000Z", true},
		{"2015-12-24T00-00-00.000Z", true},
		{"2015-12-25T00-00-00.000Z", true},
		{"2015-12-26T00-00-00.000Z", true},
		{"2015-12-27T00-00-00.000Z", true},
		{"2015-12-28T00-00-00.000Z", true},
		{"2015-12-29T00-00-00.000Z", true},
		{"2015-12-30T00-00-00.000Z", true},
		{"2015-12-31T00-00-00.000Z", true},
	}
	entries := createEntries(testDirectories, t)

	// Act
	prune := NewPrune(config)
	pruneResult, err := prune.Calculate(entries)
	if err != nil {
		t.Fatalf("Failed to calculate directories to prune: %s", err)
	}

	// Assert
	// if expected := 21; len(pruneResult.ToKeep) != expected {
	// 	t.Fatalf("Got %v, expected %v", len(pruneResult.ToKeep), expected)
	// }

	assertResultMatchesTestObjects(testDirectories, pruneResult, t)
}

// createEntries creates a list of TimeStampedDirectory based on a list of test objects
func createEntries(testObjects []TestObject, t *testing.T) []TimeStampedDirectory {
	virtualDirectories := []fs.DirEntry{}
	for _, dir := range testObjects {
		virtualDirectories = append(virtualDirectories, NewVirtualDirEntry(dir.Name, true))
	}

	entries, err := Parse(testBaseDirectory, virtualDirectories)
	if err != nil {
		t.Fatalf("Failed to parse directories: %s", err)
	}

	return entries
}

func assertResultMatchesTestObjects(testObjects []TestObject, result PruneResult, t *testing.T) {
	for _, testObject := range testObjects {
		if resultObject, ok := result.Objects[path.Join(testBaseDirectory, testObject.Name)]; ok {
			if resultObject.Keep != testObject.ExpectedKeep {
				t.Errorf("%v: Got %v, expected %v", resultObject.Directory.Time, resultObject.Keep, testObject.ExpectedKeep)
			}
		} else {
			t.Errorf("assertResultMatchesTestObjects: Expected %v to be present in result", testObject.Name)
		}
	}
}

type TestObject struct {
	Name         string
	ExpectedKeep bool
}

type VirtualDirEntry struct {
	_Name  string
	_IsDir bool
}

func NewVirtualDirEntry(name string, isDir bool) VirtualDirEntry {
	return VirtualDirEntry{_Name: name, _IsDir: isDir}
}

func (t VirtualDirEntry) Name() string {
	return t._Name
}

func (t VirtualDirEntry) IsDir() bool {
	return t._IsDir
}

func (t VirtualDirEntry) Type() fs.FileMode {
	panic("Not supported")
}

func (t VirtualDirEntry) Info() (fs.FileInfo, error) {
	panic("Not supported")
}
