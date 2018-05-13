package main

import (
	"fmt"
	"sort"
	"time"
)

type TimeRangeSubtractionResult struct {
	result              TimeRanges
	remainingSubtractor TimeRanges
}

// TimeRange represents a single TimeRange with start and end time
type TimeRange struct {
	start time.Time
	end   time.Time
}

func (tr TimeRange) Subtract(other TimeRange) TimeRangeSubtractionResult {
	result := TimeRangeSubtractionResult{}
	// s1, e1, s2, e2
	// s1, s2, e1, e2
	// s1, s2, e2, e1
	// s2, e2, s1, e1
	// s2, s1, e1, e2
	// s2, s1, e2, e1
	// s1==s2, e1, e2
	// s1==s2, e2, e1
	// s1, s2, e1==e2
	// s2, s1, e1==e2
	if tr.start.Before(other.start) && tr.end {

	}
	return result
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("(%s - %s)", tr.start.Format("15:04"), tr.end.Format("15:04"))
}

// TimeRange represents a list of TimeRanges
type TimeRanges []TimeRange

func (trs TimeRanges) Len() int {
	return len(trs)
}

func (trs TimeRanges) Less(i, j int) bool {
	if trs[i].start.Before(trs[j].start) {
		return true
	}
	if trs[i].start.After(trs[j].start) {
		return false
	}
	return trs[i].end.Before(trs[j].end)
}

func (trs TimeRanges) Swap(i, j int) {
	trs[i], trs[j] = trs[j], trs[i]
}

func (trs TimeRanges) Subtract(other TimeRanges) TimeRanges {
	sort.Sort(trs)
	sort.Sort(other)

	timeRanges := TimeRanges{}
	var i, j int // i is the index for trs. j is the index for other
	for {
		if j >= len(other) {
			break
		}

		i += 1
		j += 1
	}
	// Append the remaining ranges of trs to the return value
	for k := 0; k < len(trs); k++ {
		timeRanges = append(timeRanges, trs[k])
	}
	return timeRanges
}
