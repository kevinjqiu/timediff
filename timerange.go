package main

import (
	"fmt"
	"time"
)

// TimeRange represents a single TimeRange with start and end time
type TimeRange struct {
	start time.Time
	end   time.Time
}

func (tr TimeRange) Subtract(other TimeRange) TimeRange {
	return TimeRange{}
}

func (tr TimeRange) String() string {
	return fmt.Sprintf("(%s - %s)", tr.start.Format("15:04"), tr.end.Format("15:04"))
}

// TimeRange represents a list of TimeRanges
type TimeRanges []TimeRange

func (trs TimeRanges) Subtract(other TimeRanges) TimeRanges {
	return TimeRanges{}
}
