package main

import (
	"sort"
	"time"
)

type Configuration struct {
	Path      string
	Pattern   string
	KeepDaily int
}

func NewConfiguration(path string, keepDaily int) Configuration {
	return Configuration{Path: path, Pattern: "", KeepDaily: keepDaily}

}

type Prune struct {
	config Configuration
}

func NewPrune(c Configuration) Prune {
	return Prune{config: c}
}

func (p *Prune) Calculate(directories []TimeStampedDirectory) (PruneResult, error) {

	// Copy to new struct with keep flag
	objects := make([]PruneCandidate, 0, len(directories))
	for _, directory := range directories {
		objects = append(objects, PruneCandidate{Directory: directory})
	}

	// Group
	groupsByDay := p.groupByDay(objects)
	p.calculateKeepByDay(groupsByDay)

	// groupsByMonth := p.groupByMonth(objects)
	// p.calculateKeepByMonth(groupsByMonth)

	keep, prune := filterTimeStampedObjectByKeep(objects)

	objectsMap := make(map[string]*PruneCandidate)
	for i := 0; i < len(objects); i++ {
		object := &objects[i]
		objectsMap[object.Directory.Path] = object
	}

	result := PruneResult{Objects: objectsMap, ToKeep: keep, ToPrune: prune}

	return result, nil
}

func (p *Prune) groupByDay(objects []PruneCandidate) map[Day][]*PruneCandidate {
	groups := make(map[Day][]*PruneCandidate)

	for i := 0; i < len(objects); i++ {
		object := &objects[i]
		year, month, rawDay := object.Directory.Time.Date()
		day := Day{year, month, rawDay, time.Date(year, month, rawDay, 0, 0, 0, 0, time.UTC)}

		if value, ok := groups[day]; ok {
			groups[day] = append(value, object)
		} else {
			groups[day] = []*PruneCandidate{object}
		}
	}

	return groups
}

func (p *Prune) calculateKeepByDay(groups map[Day][]*PruneCandidate) {
	// get a sorted slice of the keys of the array
	keys := make([]Day, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]
		return a.Time.After(b.Time)
	})

	keepCount := 0
	for _, key := range keys {
		if keepCount == p.config.KeepDaily {
			break
		}

		dailyObjects := groups[key]
		var dailyToKeep *PruneCandidate

		if l := len(dailyObjects); l == 0 {
			continue
		}
		if l := len(dailyObjects); l == 1 {
			// Take single
			dailyToKeep = dailyObjects[0]
		} else if l > 1 {
			// Take oldest
			sort.SliceStable(dailyObjects, func(i, j int) bool {
				return dailyObjects[i].Directory.Time.After(dailyObjects[j].Directory.Time)
			})
			dailyToKeep = dailyObjects[0]
		}

		// Set keep if not yet set
		if !dailyToKeep.Keep {
			dailyToKeep.Keep = true
			keepCount++
		}
	}
}

func filterTimeStampedObjectByKeep(objects []PruneCandidate) ([]PruneCandidate, []PruneCandidate) {
	keep := make([]PruneCandidate, 0, len(objects))
	prune := make([]PruneCandidate, 0, len(objects))

	for _, object := range objects {
		if object.Keep {
			keep = append(keep, object)
		} else {
			prune = append(prune, object)
		}
	}

	return keep, prune
}

type PruneCandidate struct {
	Directory TimeStampedDirectory
	Keep      bool
}

type Day struct {
	Year  int
	Month time.Month
	Day   int
	Time  time.Time
}

type PruneResult struct {
	Objects map[string]*PruneCandidate
	ToKeep  []PruneCandidate
	ToPrune []PruneCandidate
}
