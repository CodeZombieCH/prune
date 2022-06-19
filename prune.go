package main

import (
	"time"
)

const NoPrune = -1

type Configuration struct {
	Path        string
	Pattern     string
	KeepDaily   int
	KeepMonthly int
	KeepYearly  int
}

func (c *Configuration) requiresPruning() bool {
	return c.KeepDaily > NoPrune && c.KeepMonthly > NoPrune && c.KeepYearly > NoPrune
}

type Prune struct {
	config Configuration
}

func NewPrune(c Configuration) Prune {
	return Prune{config: c}
}

func (p *Prune) Calculate(directories []TimeStampedDirectory) (PruneResult, error) {
	// Return immediately if empty set of directories
	if len(directories) == 0 {
		return PruneResult{Objects: make(map[string]*PruneCandidate), ToKeep: []PruneCandidate{}, ToPrune: []PruneCandidate{}}, nil
	}

	// Copy to new struct with keep flag
	objects := make([]PruneCandidate, 0, len(directories))
	for _, directory := range directories {
		objects = append(objects, PruneCandidate{Directory: directory})
	}

	if p.config.requiresPruning() {
		// Currently we do not use an array/slice, as we need the rules to be applied in a very specific order
		if p.config.KeepDaily > 0 {
			rule := KeepDailyRule{KeepCount: p.config.KeepDaily}
			rule.Apply(objects)
		}
		if p.config.KeepMonthly > 0 {
			rule := KeepMonthlyRule{KeepCount: p.config.KeepMonthly}
			rule.Apply(objects)
		}
		if p.config.KeepYearly > 0 {
			rule := KeepYearlyRule{KeepCount: p.config.KeepYearly}
			rule.Apply(objects)
		}
	} else {
		// Nothing to prune, set the keep flag on all objects
		for i := 0; i < len(objects); i++ {
			object := &objects[i]
			object.Keep = true
		}
	}

	keep, prune := filterTimeStampedObjectByKeep(objects)

	objectsMap := make(map[string]*PruneCandidate)
	for i := 0; i < len(objects); i++ {
		object := &objects[i]
		objectsMap[object.Directory.Path] = object
	}

	result := PruneResult{Objects: objectsMap, ToKeep: keep, ToPrune: prune}

	return result, nil
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
