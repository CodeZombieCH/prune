package main

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestGetObjectsDateAndTimePattern(t *testing.T) {
	rootDir := t.TempDir()

	// Arrange
	directories := []struct {
		name         string
		expectedTime time.Time
	}{
		{
			name:         "2000-12-31T00-00-00Z",
			expectedTime: time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "2001-01-01T00-00-00Z",
			expectedTime: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "2001-01-02T03-00-00Z",
			expectedTime: time.Date(2001, 1, 2, 3, 0, 0, 0, time.UTC),
		},
		{
			name:         "2001-01-02T03-04-00Z",
			expectedTime: time.Date(2001, 1, 2, 3, 4, 0, 0, time.UTC),
		},
		{
			name:         "2001-01-02T03-04-05Z",
			expectedTime: time.Date(2001, 1, 2, 3, 4, 5, 0, time.UTC),
		},
	}

	for _, v := range directories {
		err := os.Mkdir(path.Join(rootDir, v.name), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s", v.name)
		}
	}

	// Act
	traverser := FileSystemTraverser{Pattern: PatternAlmostISO8601DateAndTime}
	objects, err := traverser.GetObjects(rootDir)

	// Assert
	if err != nil {
		t.Fatalf("Failed to get objects for path %s: %v", rootDir, err)
	}

	if expected, actual := 5, len(objects); expected != actual {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}

	objectsMap := toObjectsMap(objects)
	for _, v := range directories {
		object, ok := objectsMap[path.Join(rootDir, v.name)]
		if !ok {
			t.Errorf("Expected %v to be present", v)
		}
		if v.expectedTime != object.Time {
			t.Fatalf("Expected %v, got %v", v.expectedTime, object.Time)
		}
	}
}

func TestGetObjectsDateOnlyPattern(t *testing.T) {
	rootDir := t.TempDir()

	// Arrange
	directories := []struct {
		name         string
		expectedTime time.Time
	}{
		{
			name:         "2000-12-30",
			expectedTime: time.Date(2000, 12, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "2000-12-31",
			expectedTime: time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "2001-01-01",
			expectedTime: time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:         "2001-01-02",
			expectedTime: time.Date(2001, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, v := range directories {
		err := os.Mkdir(path.Join(rootDir, v.name), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s", v.name)
		}
	}

	// Act
	traverser := FileSystemTraverser{Pattern: PatternISO8601DateOnly}
	objects, err := traverser.GetObjects(rootDir)

	// Assert
	if err != nil {
		t.Fatalf("Failed to get objects for path %s: %v", rootDir, err)
	}

	if expected, actual := 4, len(objects); expected != actual {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}

	objectsMap := toObjectsMap(objects)
	for _, v := range directories {
		object, ok := objectsMap[path.Join(rootDir, v.name)]
		if !ok {
			t.Errorf("Expected %v to be present", v)
		}
		if v.expectedTime != object.Time {
			t.Fatalf("Expected %v, got %v", v.expectedTime, object.Time)
		}
	}
}

// Expected to fail, but date is automatically "corrected"
// 2000-02-31 => 2000-03-02 00:00:00 +0000 UTC
// 2000-13-31 => 2001-01-31 00:00:00 +0000 UTC
// Not sure if I like this or not
func TestGetObjectsInvalidDates(t *testing.T) {
	rootDir := t.TempDir()

	// Arrange
	directories := []struct {
		name         string
		expectedTime time.Time
	}{
		{
			name: "2000-02-31",
		},
		{
			name: "2000-13-31",
		},
	}

	for _, v := range directories {
		err := os.Mkdir(path.Join(rootDir, v.name), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s", v.name)
		}
	}

	// Act
	traverser := FileSystemTraverser{Pattern: PatternISO8601DateOnly}
	objects, err := traverser.GetObjects(rootDir)

	if err != nil {
		t.Fatalf("Failed to get objects for path %s: %v", rootDir, err)
	}

	if expected, actual := 2, len(objects); expected != actual {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func TestGetObjectsInvalidPattern(t *testing.T) {
	rootDir := t.TempDir()

	// Arrange
	directories := []struct {
		name         string
		expectedTime time.Time
	}{
		{
			name: "02.01.2000",
		},
		{
			name: "foo-2000-12-31-bar",
		},
	}

	for _, v := range directories {
		err := os.Mkdir(path.Join(rootDir, v.name), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s", v.name)
		}
	}

	// Act
	traverser := FileSystemTraverser{Pattern: PatternISO8601DateOnly}
	objects, err := traverser.GetObjects(rootDir)

	if err != nil {
		t.Fatalf("Failed to get objects for path %s: %v", rootDir, err)
	}

	if expected, actual := 0, len(objects); expected != actual {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func toObjectsMap(objects []TimeStampedDirectory) map[string]*TimeStampedDirectory {
	objectsMap := make(map[string]*TimeStampedDirectory)

	for i := 0; i < len(objects); i++ {
		object := &objects[i]
		objectsMap[object.Path] = object
	}

	return objectsMap
}
