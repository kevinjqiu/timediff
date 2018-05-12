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
func makeTimeRangeHHMM(start string, end string) TimeRange {
	return TimeRange{
		makeHHMM(start),
		makeHHMM(end),
	}
}

func TestTimeRange1SupersedesTimeRange2(t *testing.T) {
	tr1 := TimeRanges{makeTimeRangeHHMM("09:00", "10:00")}
	tr2 := TimeRanges{makeTimeRangeHHMM("09:00", "09:30")}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{makeTimeRangeHHMM("09:30", "10:00")})
}

func TestTimeRange1EqualsTimeRange2(t *testing.T) {
	tr1 := TimeRanges{makeTimeRangeHHMM("09:30", "10:30")}
	tr2 := TimeRanges{makeTimeRangeHHMM("09:30", "10:30")}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{})
}

func TestTimeRange1BoardersTimeRange2(t *testing.T) {
	tr1 := TimeRanges{makeTimeRangeHHMM("09:00", "09:30")}
	tr2 := TimeRanges{makeTimeRangeHHMM("09:30", "15:00")}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{makeTimeRangeHHMM("09:00", "09:30")})
}

func TestTimeRangesDoNotIntersect(t *testing.T) {
	tr1 := TimeRanges{makeTimeRangeHHMM("09:00", "09:30")}
	tr2 := TimeRanges{makeTimeRangeHHMM("09:31", "15:00")}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{makeTimeRangeHHMM("09:00", "09:30")})
}

func TestMultiTimeRangeOverlap(t *testing.T) {
	tr1 := TimeRanges{
		makeTimeRangeHHMM("09:00", "09:30"),
		makeTimeRangeHHMM("10:00", "10:30"),
	}
	tr2 := TimeRanges{
		makeTimeRangeHHMM("09:15", "10:15"),
	}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{
		makeTimeRangeHHMM("09:00", "09:15"),
		makeTimeRangeHHMM("10:15", "10:30"),
	})
}

func TestMultiTimeRangeSubtraction(t *testing.T) {
	tr1 := TimeRanges{
		makeTimeRangeHHMM("09:00", "11:00"),
		makeTimeRangeHHMM("13:00", "15:00"),
	}
	tr2 := TimeRanges{
		makeTimeRangeHHMM("09:00", "09:15"),
		makeTimeRangeHHMM("10:00", "10:15"),
		makeTimeRangeHHMM("12:30", "16:00"),
	}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{
		makeTimeRangeHHMM("09:15", "10:00"),
		makeTimeRangeHHMM("10:15", "11:00"),
	})
}

func TestMultiTimeRangeSubtractionOutOfOrder(t *testing.T) {
	tr1 := TimeRanges{
		makeTimeRangeHHMM("13:00", "15:00"),
		makeTimeRangeHHMM("09:00", "11:00"),
	}
	tr2 := TimeRanges{
		makeTimeRangeHHMM("10:00", "10:15"),
		makeTimeRangeHHMM("09:00", "09:15"),
		makeTimeRangeHHMM("12:30", "16:00"),
	}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{
		makeTimeRangeHHMM("09:15", "10:00"),
		makeTimeRangeHHMM("10:15", "11:00"),
	})
}
