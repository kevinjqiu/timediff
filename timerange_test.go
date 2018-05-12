package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func makeHHMM(timeStr string) time.Time {
	timeStruct, err := time.Parse("15:04", timeStr)
	if err != nil {
		panic(err)
	}
	return timeStruct
}

// Construct a TimeRanges object
func makeTimeRangesHHMM(start string, end string) TimeRanges {
	return TimeRanges{
		TimeRange{
			makeHHMM(start),
			makeHHMM(end),
		},
	}
}

func TestTimeRange1SupersedesTimeRange2(t *testing.T) {
	tr1 := makeTimeRangesHHMM("09:00", "10:00")
	tr2 := makeTimeRangesHHMM("09:00", "09:30")
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, makeTimeRangesHHMM("09:30", "10:00"))
}

func TestTimeRange1EqualsTimeRange2(t *testing.T) {
	tr1 := makeTimeRangesHHMM("09:30", "10:30")
	tr2 := makeTimeRangesHHMM("09:30", "10:30")
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRange{})
}

func TestTimeRange1BoardersTimeRange2(t *testing.T) {
	tr1 := makeTimeRangesHHMM("09:00", "09:30")
	tr2 := makeTimeRangesHHMM("09:30", "15:00")
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, makeTimeRangesHHMM("09:00", "09:30"))
}

func TestTimeRangesDoNotIntersect(t *testing.T) {
	tr1 := makeTimeRangesHHMM("09:00", "09:30")
	tr2 := makeTimeRangesHHMM("09:31", "15:00")
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, makeTimeRangesHHMM("09:00", "09:30"))
}
