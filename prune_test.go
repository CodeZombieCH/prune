package main

import (
	"io/fs"
	"path"
	"testing"
)

var baseDirectory string = "/foo/bar"

func TestPruneWithMultiplePerDay(t *testing.T) {
	// Arrange
	config := NewConfiguration(baseDirectory, 4)
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
	config := NewConfiguration(baseDirectory, 7)
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

// createEntries creates a list of TimeStampedDirectory based on a list of test objects
func createEntries(testObjects []TestObject, t *testing.T) []TimeStampedDirectory {
	virtualDirectories := []fs.DirEntry{}
	for _, dir := range testObjects {
		virtualDirectories = append(virtualDirectories, NewVirtualDirEntry(dir.Name, true))
	}

	entries, err := Parse(baseDirectory, virtualDirectories)
	if err != nil {
		t.Fatalf("Failed to parse directories: %s", err)
	}

	return entries
}

func assertResultMatchesTestObjects(testObjects []TestObject, result PruneResult, t *testing.T) {
	for _, testObject := range testObjects {
		if resultObject, ok := result.Objects[path.Join(baseDirectory, testObject.Name)]; ok {
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
