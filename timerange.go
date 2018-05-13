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

func (trsr TimeRangeSubtractionResult) String() string {
	return fmt.Sprintf("result=%v, remaining=%v", trsr.result, trsr.remainingSubtractor)
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
func (tr TimeRange) Subtract(subtractor TimeRange) TimeRangeSubtractionResult {
	result := TimeRangeSubtractionResult{
		result:              TimeRanges{},
		remainingSubtractor: TimeRanges{},
	}
	s1 := tr.start
	s2 := subtractor.start
	e1 := tr.end
	e2 := subtractor.end

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

// ReplaceAt replaces the element at idx with elements in the newTrs
func (trs TimeRanges) ReplaceAt(idx int, newTrs TimeRanges) TimeRanges {
	result := append(trs[0:idx], newTrs...)
	if idx <= len(trs)-1 {
		result = append(result, trs[idx+1:len(trs)]...)
	}
	return result
}

// IsEmpty returns true if the TimeRanges slice is empty
func (trs TimeRanges) IsEmpty() bool {
	return trs.Len() == 0
}

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

// Merge all the time ranges in this TimeRanges object
// Precondition: the time ranges are sorted
func (trs TimeRanges) Merge() TimeRanges {
	if trs.Len() < 2 {
		return trs
	}

	newTrs := TimeRanges{}
	mergingTr := trs[0] // Keep a reference to the TimeRange currently being merged

	for _, tr := range trs[1:] {
		switch {
		case tr.start.After(mergingTr.end):
			newTrs = append(newTrs, mergingTr)
			mergingTr = tr
		case (tr.start.Before(mergingTr.end) || tr.start.Equal(mergingTr.end)) && (tr.end.After(mergingTr.end) || tr.end.Equal(mergingTr.end)):
			mergingTr.end = tr.end
		}
	}

	// Add the remaining
	newTrs = append(newTrs, mergingTr)
	return newTrs
}

// Subtract two TimeRanges object, returns a new TimeRanges object containing the result
func (trs TimeRanges) Subtract(subtractors TimeRanges) TimeRanges {
	sort.Sort(trs)
	sort.Sort(subtractors)

	var (
		i, j int
		diff TimeRangeSubtractionResult
	)
	for {
		if i >= len(trs) {
			break
		}
		if j >= len(subtractors) {
			break
		}

		tr := trs[i]
		subtractor := subtractors[j]
		diff = tr.Subtract(subtractor)
		// fmt.Println("##################")
		// fmt.Printf("i=%d\n", i)
		// fmt.Printf("j=%d\n", j)
		// fmt.Printf("trs=%v\n", trs)
		// fmt.Printf("subtractors=%v\n", subtractors)
		// fmt.Printf("tr=%v\n", tr)
		// fmt.Printf("subtractor=%v\n", subtractor)
		// fmt.Printf("diff=%v\n", diff)
		// fmt.Println("-----------------")

		trs = trs.ReplaceAt(i, diff.result)
		if len(diff.result) >= 1 && diff.result[0] == tr {
			i++
			continue
		}
		subtractors = subtractors.ReplaceAt(j, diff.remainingSubtractor)
		if len(diff.remainingSubtractor) == 1 && diff.remainingSubtractor[0] == subtractor {
			j++
			continue
		}

		// fmt.Printf("i=%d\n", i)
		// fmt.Printf("j=%d\n", j)
		// fmt.Println("-----------------")
		// fmt.Printf("trs=%v\n", trs)
		// fmt.Printf("subtractors=%v\n", subtractors)
		// fmt.Println("##################")
	}
	return trs.Merge()
}
