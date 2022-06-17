package main

import (
	"os"
	"path"
	"testing"
)

func TestGetObjects(t *testing.T) {
	rootDir := t.TempDir()

	// Arrange
	dirs := []string{
		"2022-06-16T01-00-00.000Z",
		"2022-06-16T02-00-00.000Z",
		"2022-06-17T01-00-00.000Z",
		"2022-06-20T01-00-00.000Z",
		"2022-06-21T01-00-00.000Z",
		"2022-06-22T01-00-00.000Z",
		"2022-06-23T01-00-00.000Z",
		"2022-06-24T01-00-00.000Z",
	}

	for _, v := range dirs {
		err := os.Mkdir(path.Join(rootDir, v), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s", v)
		}
	}

	// Act
	traverser := FileSystemTraverser{}
	objects, err := traverser.GetObjects(rootDir)

	// Assert
	if err != nil {
		t.Fatalf("Failed to get objects for path: %s", rootDir)
	}

	if len(objects) != 8 {
		t.Fatalf("Got %v, expected %v", len(objects), 8)
	}
}
