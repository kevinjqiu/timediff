package main

import (
	"fmt"
	"sort"
	"time"
)

// TimeRangeSubtractionResult represents the result of `TimeRange.Subtract` method
type TimeRangeSubtractionResult struct {
	result              TimeRanges // The result of applying tr1-tr2
	remainingSubtractor TimeRanges // The remaining of tr2 after applying tr1-tr2 that's after the start of tr1
}

// TimeRange represents a single TimeRange with start and end time
type TimeRange struct {
	start time.Time
	end   time.Time
}

/*
Subtract one TimeRange from another
Let's denote `tr1(start=s1, end=e1)` and `tr2(start=e2, end=e2)` be the two TimeRange objects.
`tr1.Subtract(tr2)` can have the following scenarios:
# Scenario 1 (s1==s2 && e1==e2)
	t (time) --------------------------------------->
	tr1             s1------------e1
	tr2             s2------------e2
	result          (empty)
	remaining       (empty)
# Scenario 2 (s1==s2 && e1<e2)
	t (time) --------------------------------------->
	tr1             s1--------e1
	tr2             s2------------e2
	result          (empty)
	remaining                 e1--e2
# Scenario 3 (s1==s2 && e1>e2)
	t (time) --------------------------------------->
	tr1             s1-----------------e1
	tr2             s2------------e2
	result                        e2---e1
	remaining       (empty)
# Scenario 4 (s1<s2 && e1==e2)
	t (time) --------------------------------------->
	tr1             s1-----------------e1
	tr2                 s2-------------e2
	result          s1--s2
	remaining       (empty)
# Scenario 5 (s1>s2 && e1==e2)
	t (time) --------------------------------------->
	tr1                s1------------------e1
	tr2             s2---------------------e2
	result          (empty)
	remaining       (empty)
# Scenario 6 (s1<s2 && e1<e2 && s2<e1)
	t (time) --------------------------------------->
	tr1             s1------------------e1
	tr2                  s2---------------------e2
	result          s1---s2
	remaining                           e1------e2
# Scenario 7 (s1<s2 && e1>e2 && s2<e1)
	t (time) --------------------------------------->
	tr1             s1--------------------------e1
	tr2                  s2---------------e2
	result          s1---s2               e2----e1
	remaining       (empty)
# Scenario 8 (s1>s2 && e1<e2 && s1<e2)
	t (time) --------------------------------------->
	tr1                  s1------e1
	tr2             s2---------------e2
	result          (empty)
	remaining                    e1--e2
# Scenario 9 (s1>s2 && e1>e2 && s1<e2)
	t (time) --------------------------------------->
	tr1                  s1-----------------e1
	tr2             s2---------------e2
	result                           e2-----e1
	remaining       (empty)
# Scenario 10 (s1>=e2)
	t (time) --------------------------------------->
	tr1                       s1-----------------e1
	tr2             s2-----e2
	result                    s1-----------------e1
	remaining       (empty)
# Scenario 11 (s2>=e1)
	t (time) --------------------------------------->
	tr1             s1-----------------e1
	tr2                                    s2-----e2
	result          s1-----------------e1
	remaining                              s2-----e2
*/
func (tr TimeRange) Subtract(other TimeRange) TimeRangeSubtractionResult {
	result := TimeRangeSubtractionResult{
		result:              TimeRanges{},
		remainingSubtractor: TimeRanges{},
	}
	s1 := tr.start
	s2 := other.start
	e1 := tr.end
	e2 := other.end

	switch {
	case s1.Equal(s2) && e1.Equal(e2):
		// Do nothing, as the default for result and remaining are empty TimeRanges slice
	case s1.Equal(s2) && e1.Before(e2):
		result.remainingSubtractor = TimeRanges{TimeRange{e1, e2}}
	case s1.Equal(s2) && e1.After(e2):
		result.result = TimeRanges{TimeRange{e2, e1}}
	case s1.Before(s2) && e1.Equal(e2):
		result.result = TimeRanges{TimeRange{s1, s2}}
	case s1.After(s2) && e1.Equal(e2):
		// Do nothing, as the default for result and remaining are empty TimeRanges slice
	case s1.Before(s2) && e1.Before(e2) && s2.Before(e1):
		result.result = TimeRanges{TimeRange{s1, s2}}
		result.remainingSubtractor = TimeRanges{TimeRange{e1, e2}}
	case s1.Before(s2) && e1.After(e2) && s2.Before(e1):
		result.result = TimeRanges{
			TimeRange{s1, s2},
			TimeRange{e2, e1},
		}
	case s1.After(s2) && e1.Before(e2) && s1.Before(e2):
		result.remainingSubtractor = TimeRanges{TimeRange{e1, e2}}
	case s1.After(s2) && e1.After(e2) && s1.Before(e2):
		result.result = TimeRanges{TimeRange{e2, e1}}
	case s1.After(e2) || s1.Equal(e2):
		result.result = TimeRanges{TimeRange{s1, e1}}
	case s2.After(e1) || s2.Equal(e1):
		result.result = TimeRanges{TimeRange{s1, e1}}
		result.remainingSubtractor = TimeRanges{TimeRange{s2, e2}}
	default:
		panic("Shouldn't reach here. Were there missing scenarios?")
	}
	return result
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("(%s - %s)", tr.start.Format("15:04"), tr.end.Format("15:04"))
}

// TimeRanges represents a list of TimeRanges
type TimeRanges []TimeRange

// Method to satisfy the Sort interface
func (trs TimeRanges) Len() int {
	return len(trs)
}

// Method to satisfy the Sort interface
func (trs TimeRanges) Less(i, j int) bool {
	if trs[i].start.Before(trs[j].start) {
		return true
	}
	if trs[i].start.After(trs[j].start) {
		return false
	}
	return trs[i].end.Before(trs[j].end)
}

// Method to satisfy the Sort interface
func (trs TimeRanges) Swap(i, j int) {
	trs[i], trs[j] = trs[j], trs[i]
}

// Subtract two TimeRanges object
func (trs TimeRanges) Subtract(other TimeRanges) TimeRanges {
	sort.Sort(trs)
	sort.Sort(other)

	timeRanges := TimeRanges{}
	var i, j int // i is the index for trs. j is the index for other
	for {
		if j >= len(other) {
			break
		}

		i++
		j++
	}
	// Append the remaining ranges of trs to the return value
	for k := 0; k < len(trs); k++ {
		timeRanges = append(timeRanges, trs[k])
	}
	return timeRanges
}
