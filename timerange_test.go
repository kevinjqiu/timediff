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

// Construct a TimeRange object using HHMM format string
func mktr(start string, end string) TimeRange {
	return TimeRange{
		makeHHMM(start),
		makeHHMM(end),
	}
}

// Tests for TimeRange.Subtract
// TR1:                |-------------|
// TR2:                |-------------|
// Result:             []
// TR2 Remaining:      []
func TestTimeRangeSubtractionSameRange(t *testing.T) {
	tr1 := mktr("09:00", "10:00")
	tr2 := mktr("09:00", "10:00")
	diff := tr1.Subtract(tr2)
	assert.Equal(t, diff.result, TimeRanges{})
	assert.Equal(t, diff.remainingSubtractor, TimeRanges{})
}

// TR1:                |-------------|
// TR2:            |----------------------|
// Result:         []
// TR2 Remaining:                    |----|
func TestTimeRangeSubtractionTR2EncompassesTR1(t *testing.T) {
	tr1 := mktr("09:00", "10:00")
	tr2 := mktr("08:00", "11:00")
	diff := tr1.Subtract(tr2)
	assert.Equal(t, diff.result, TimeRanges{})
	assert.Equal(t, diff.remainingSubtractor, TimeRanges{
		mktr("09:00", "11:00"),
	})
}

// TR1:            |----------------------|
// TR2:                |-------------|
// Result:         |---|             |----|
// TR2 Remaining:  []
func TestTimeRangeSubtractionTR2BisectsTR1(t *testing.T) {
	tr1 := mktr("08:00", "11:00")
	tr2 := mktr("09:00", "10:00")
	diff := tr1.Subtract(tr2)
	assert.Equal(t, diff.result, TimeRanges{
		mktr("08:00", "09:00"),
		mktr("09:00", "11:00"),
	})
	assert.Equal(t, diff.remainingSubtractor, TimeRanges{})
}

// TR1:            |-----|
// TR2:                    |-------------|
// Result:         |-----|
// TR2 Remaining:          |-------------|
func TestTimeRangeSubtractionTR2DoesNotIntersectTR1(t *testing.T) {
	tr1 := mktr("08:00", "11:00")
	tr2 := mktr("12:00", "15:00")
	diff := tr1.Subtract(tr2)
	assert.Equal(t, diff.result, TimeRanges{tr1})
	assert.Equal(t, diff.remainingSubtractor, TimeRanges{tr2})
}

// TR1:            |--------|
// TR2:                |-------------|
// Result:         |---|
// TR2 Remaining:           |--------|
func TestTimeRangeSubtractionTR2OverlapsAndIsLaterThanTR1(t *testing.T) {
	tr1 := mktr("08:00", "11:00")
	tr2 := mktr("10:00", "15:00")
	diff := tr1.Subtract(tr2)
	assert.Equal(t, diff.result, TimeRange{mktr("08:00", "10:00")})
	assert.Equal(t, diff.remainingSubtractor, TimeRange{mktr("11:00", "15:00")})
}

// TR1:                    |--------|
// TR2:            |-------------|
// Result:                       |--|
// TR2 Remaining:  []
func TestTimeRangeSubtractionTR2OverlapsAndIsEarlierThanTR1(t *testing.T) {
	tr1 := mktr("08:00", "11:00")
	tr2 := mktr("07:00", "09:00")
	diff := tr1.Subtract(tr2)
	assert.Equal(t, diff.result, TimeRanges{mktr("09:00", "11:00")})
	assert.Equal(t, diff.remainingSubtractor, TimeRanges{})
}

// Tests for TimeRanges.Subtract
func TestTimeRange1SupersedesTimeRange2(t *testing.T) {
	tr1 := TimeRanges{mktr("09:00", "10:00")}
	tr2 := TimeRanges{mktr("09:00", "09:30")}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{mktr("09:30", "10:00")})
}

func TestTimeRange1EqualsTimeRange2(t *testing.T) {
	tr1 := TimeRanges{mktr("09:30", "10:30")}
	tr2 := TimeRanges{mktr("09:30", "10:30")}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{})
}

func TestTimeRange1BoardersTimeRange2(t *testing.T) {
	tr1 := TimeRanges{mktr("09:00", "09:30")}
	tr2 := TimeRanges{mktr("09:30", "15:00")}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{mktr("09:00", "09:30")})
}

func TestTimeRangesDoNotIntersect(t *testing.T) {
	tr1 := TimeRanges{mktr("09:00", "09:30")}
	tr2 := TimeRanges{mktr("09:31", "15:00")}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{mktr("09:00", "09:30")})
}

func TestMultiTimeRangeOverlap(t *testing.T) {
	tr1 := TimeRanges{
		mktr("09:00", "09:30"),
		mktr("10:00", "10:30"),
	}
	tr2 := TimeRanges{
		mktr("09:15", "10:15"),
	}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{
		mktr("09:00", "09:15"),
		mktr("10:15", "10:30"),
	})
}

func TestMultiTimeRangeSubtraction(t *testing.T) {
	tr1 := TimeRanges{
		mktr("09:00", "11:00"),
		mktr("13:00", "15:00"),
	}
	tr2 := TimeRanges{
		mktr("09:00", "09:15"),
		mktr("10:00", "10:15"),
		mktr("12:30", "16:00"),
	}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{
		mktr("09:15", "10:00"),
		mktr("10:15", "11:00"),
	})
}

func TestMultiTimeRangeSubtractionOutOfOrder(t *testing.T) {
	tr1 := TimeRanges{
		mktr("13:00", "15:00"),
		mktr("09:00", "11:00"),
	}
	tr2 := TimeRanges{
		mktr("10:00", "10:15"),
		mktr("09:00", "09:15"),
		mktr("12:30", "16:00"),
	}
	trdiff := tr1.Subtract(tr2)
	assert.Equal(t, trdiff, TimeRanges{
		mktr("09:15", "10:00"),
		mktr("10:15", "11:00"),
	})
}
