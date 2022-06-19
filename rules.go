package main

import (
	"sort"
	"time"
)

type Rule interface {
	Apply(objects []PruneCandidate)
}

type KeepDailyRule struct {
	KeepCount int
}

func (r *KeepDailyRule) Apply(objects []PruneCandidate) {
	groups := groupBy(objects, KeepDailyTimeConvert)
	applyKeepRule(groups, r.KeepCount)
}

func KeepDailyTimeConvert(exactTime time.Time) time.Time {
	year, month, day := exactTime.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

type KeepMonthlyRule struct {
	KeepCount int
}

func (r *KeepMonthlyRule) Apply(objects []PruneCandidate) {
	groups := groupBy(objects, KeepMonthlyTimeConvert)
	applyKeepRule(groups, r.KeepCount)
}

func KeepMonthlyTimeConvert(exactTime time.Time) time.Time {
	return time.Date(exactTime.Year(), exactTime.Month(), 0, 0, 0, 0, 0, time.UTC)
}

type KeepYearlyRule struct {
	KeepCount int
}

func (r *KeepYearlyRule) Apply(objects []PruneCandidate) {
	groups := groupBy(objects, KeepYearlyTimeConvert)
	applyKeepRule(groups, r.KeepCount)
}

func KeepYearlyTimeConvert(exactTime time.Time) time.Time {
	return time.Date(exactTime.Year(), 0, 0, 0, 0, 0, 0, time.UTC)
}

func groupBy(objects []PruneCandidate, timeConvert func(time time.Time) time.Time) map[time.Time][]*PruneCandidate {
	groups := make(map[time.Time][]*PruneCandidate)

	for i := 0; i < len(objects); i++ {
		object := &objects[i]
		relevantTime := timeConvert(object.Directory.Time)

		if value, ok := groups[relevantTime]; ok {
			groups[relevantTime] = append(value, object)
		} else {
			groups[relevantTime] = []*PruneCandidate{object}
		}
	}

	return groups
}

func applyKeepRule(groups map[time.Time][]*PruneCandidate, keepCount int) {
	// get a sorted slice of the keys of the array
	keys := make([]time.Time, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		// WARNING: not a less function, but a more function, so we can start from the beginning of the slice
		return keys[i].After(keys[j])
	})

	currentKeepCount := 0
	for _, key := range keys {
		if currentKeepCount == keepCount {
			break
		}

		relevantTimeObjects := groups[key]

		objectToKeep := getNewest(relevantTimeObjects)

		if objectToKeep == nil {
			continue
		}

		// Set keep if not yet set
		if !objectToKeep.Keep {
			objectToKeep.Keep = true
			currentKeepCount++
		}
	}

	// Keep oldest object
	if currentKeepCount < keepCount {
		latestKey := keys[len(keys)-1]
		relevantTimeObjects := groups[latestKey]

		objectToKeep := getOldest(relevantTimeObjects)

		if objectToKeep != nil {
			// Set keep if not yet set
			if !objectToKeep.Keep {
				objectToKeep.Keep = true
				currentKeepCount++
			}
		}
	}
}

func getNewest(candidates []*PruneCandidate) *PruneCandidate {
	return firstOrNilSorted(candidates, sortAndTakeNewest)
}

func getOldest(candidates []*PruneCandidate) *PruneCandidate {
	return firstOrNilSorted(candidates, sortAndTakeOldest)
}

func firstOrNilSorted(candidates []*PruneCandidate, sortAndFirstFunc func(candidates []*PruneCandidate) *PruneCandidate) *PruneCandidate {
	l := len(candidates)
	if l == 0 {
		return nil
	} else if l == 1 {
		// Take single
		return candidates[0]
	} else {
		return sortAndFirstFunc(candidates)
	}
}

func sortAndTakeNewest(candidates []*PruneCandidate) *PruneCandidate {
	// Take newest
	sort.SliceStable(candidates, func(i, j int) bool {
		// WARNING: not a less function, but a more function, so we can take the first element
		return candidates[i].Directory.Time.After(candidates[j].Directory.Time)
	})
	return candidates[0]
}

func sortAndTakeOldest(candidates []*PruneCandidate) *PruneCandidate {
	// Take oldest
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Directory.Time.Before(candidates[j].Directory.Time)
	})
	return candidates[0]
}
