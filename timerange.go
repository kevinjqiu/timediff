package main

import "time"

type TimeRange struct {
	start time.Time
	end   time.Time
}

func (tr TimeRange) Subtract(other TimeRange) TimeRange {
	return TimeRange{}
}
