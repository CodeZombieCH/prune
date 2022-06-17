package main

import (
	"io/fs"
	"log"
	"os"
	"path"
	"time"
)

type FileSystemTraverser struct {
}

func (t *FileSystemTraverser) GetObjects(basePath string) ([]TimeStampedDirectory, error) {
	// TODO: think about using File.Readdirnames as it should be much faster
	// TODO: think about using a channel to send file names to for further processing

	// CHECK https://bitfieldconsulting.com/golang/filesystems for more inspiration

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	objects, err := Parse(basePath, entries)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func Parse(basePath string, entries []fs.DirEntry) ([]TimeStampedDirectory, error) {

	objects := []TimeStampedDirectory{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Parse
		name := entry.Name() // Read once and cache to reduce system calls
		layout := "2006-01-02T15-04-05.000Z"
		t, err := time.Parse(layout, name)
		if err != nil {
			log.Printf("getObjects: failed to parse date for directory entry %v: %v", entry, err)
			continue
		}

		objects = append(objects, TimeStampedDirectory{Name: name, Path: path.Join(basePath, name), Time: t})
	}

	return objects, nil
}

type TimeStampedDirectory struct {
	Name string
	Path string
	Time time.Time
}
