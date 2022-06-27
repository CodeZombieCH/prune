package main

import (
	"io/fs"
	"log"
	"os"
	"path"
	"time"

	"github.com/itchyny/timefmt-go"
)

const PatternISO8601DateOnly = "%Y-%m-%d"
const PatternAlmostISO8601DateAndTime = "%Y-%m-%dT%H-%M-%S%z"

type FileSystemTraverser struct {
	Pattern string
}

func (t *FileSystemTraverser) GetObjects(basePath string) ([]TimeStampedDirectory, error) {
	// TODO: think about using File.Readdirnames as it should be much faster
	// TODO: think about using a channel to send file names to for further processing

	// CHECK https://bitfieldconsulting.com/golang/filesystems for more inspiration

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	objects, err := Parse(basePath, t.Pattern, entries)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func Parse(basePath string, pattern string, entries []fs.DirEntry) ([]TimeStampedDirectory, error) {

	objects := []TimeStampedDirectory{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Parse
		name := entry.Name() // Read once and cache to reduce system calls

		t, err := timefmt.Parse(name, pattern)
		if err != nil {
			log.Printf("getObjects: failed to parse date for directory entry %v: %v", name, err)
			continue
		}

		objects = append(objects, TimeStampedDirectory{Name: name, Path: path.Join(basePath, name), Time: t})
	}

	// Issue warning when no directory was matched by the pattern
	// TODO: should we return an error?
	if len(entries) > 0 && len(objects) == 0 {
		log.Printf("traverse: failed to parse date for all directory entries. Is your pattern '%v' valid?", pattern)
	}

	return objects, nil
}

type TimeStampedDirectory struct {
	Name string
	Path string
	Time time.Time
}
