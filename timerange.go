package timerangesubtraction

import (
	"fmt"
	"sort"
	"time"
)

// TimeRangeSubtractionResult represents the result of `TimeRange.Subtract` method
type TimeRangeSubtractionResult struct {
	result              TimeRanges // The result of applying tr1-tr2
	remainingSubtractor TimeRange  // The remaining of tr2 after applying tr1-tr2 that's after the start of tr1
}

// HasRemainingSubtractor returns true if the remainingSubtractor is empty
func (trsr TimeRangeSubtractionResult) HasRemainingSubtractor() bool {
	return trsr.remainingSubtractor != TimeRange{}
}

func (trsr TimeRangeSubtractionResult) String() string {
	return fmt.Sprintf("{result=%v, remaining=%v}", trsr.result, trsr.remainingSubtractor)
}

// TimeRange represents a single TimeRange with start and end time
type TimeRange struct {
	Start time.Time
	End   time.Time
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
		remainingSubtractor: TimeRange{},
	}
	s1 := tr.Start
	s2 := subtractor.Start
	e1 := tr.End
	e2 := subtractor.End

	switch {
	case s1.Equal(s2) && e1.Equal(e2):
		// Do nothing, as the default for result and remaining are empty TimeRanges slice
	case s1.Equal(s2) && e1.Before(e2):
		result.remainingSubtractor = TimeRange{e1, e2}
	case s1.Equal(s2) && e1.After(e2):
		result.result = TimeRanges{TimeRange{e2, e1}}
	case s1.Before(s2) && e1.Equal(e2):
		result.result = TimeRanges{TimeRange{s1, s2}}
	case s1.After(s2) && e1.Equal(e2):
		// Do nothing, as the default for result and remaining are empty TimeRanges slice
	case s1.Before(s2) && e1.Before(e2) && s2.Before(e1):
		result.result = TimeRanges{TimeRange{s1, s2}}
		result.remainingSubtractor = TimeRange{e1, e2}
	case s1.Before(s2) && e1.After(e2) && s2.Before(e1):
		result.result = TimeRanges{
			TimeRange{s1, s2},
			TimeRange{e2, e1},
		}
	case s1.After(s2) && e1.Before(e2) && s1.Before(e2):
		result.remainingSubtractor = TimeRange{e1, e2}
	case s1.After(s2) && e1.After(e2) && s1.Before(e2):
		result.result = TimeRanges{TimeRange{e2, e1}}
	case s1.After(e2) || s1.Equal(e2):
		result.result = TimeRanges{TimeRange{s1, e1}}
	case s2.After(e1) || s2.Equal(e1):
		result.result = TimeRanges{TimeRange{s1, e1}}
		result.remainingSubtractor = TimeRange{s2, e2}
	default:
		panic("Shouldn't reach here. Were there missing scenarios?")
	}
	return result
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("(%s - %s)", tr.Start.Format("15:04"), tr.End.Format("15:04"))
}

// TimeRanges represents a list of TimeRanges
type TimeRanges []TimeRange

// ReplaceAt replaces the element at idx with elements in the newTrs
func (trs TimeRanges) ReplaceAt(idx int, newTrs TimeRanges) TimeRanges {
	empty := TimeRange{}
	result := trs[0:idx]
	for _, tr := range newTrs {
		if tr != empty {
			result = append(result, tr)
		}
	}
	if idx <= len(trs)-1 {
		result = append(result, trs[idx+1:len(trs)]...)
	}
	return result
}

// Method to satisfy the Sort interface
func (trs TimeRanges) Len() int {
	return len(trs)
}

// Method to satisfy the Sort interface
func (trs TimeRanges) Less(i, j int) bool {
	if trs[i].Start.Before(trs[j].Start) {
		return true
	}
	if trs[i].Start.After(trs[j].Start) {
		return false
	}
	return trs[i].End.Before(trs[j].End)
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
		case tr.Start.After(mergingTr.End):
			newTrs = append(newTrs, mergingTr)
			mergingTr = tr
		case (tr.Start.Before(mergingTr.End) || tr.Start.Equal(mergingTr.End)) && (tr.End.After(mergingTr.End) || tr.End.Equal(mergingTr.End)):
			mergingTr.End = tr.End
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

		trs = trs.ReplaceAt(i, diff.result)
		if !diff.HasRemainingSubtractor() {
			// If the subtractor is totally consumed, move the subtractor pointer
			// and keep the subtractee (trs) pointer untouched
			j++
			continue
		}
		subtractors = subtractors.ReplaceAt(j, TimeRanges{diff.remainingSubtractor})
		if len(diff.result) >= 1 && diff.result[0] == tr {
			i++
		}
	}
	return trs.Merge()
}
