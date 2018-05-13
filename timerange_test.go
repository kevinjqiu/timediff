package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

type timeRangeSubtractTestCase struct {
	description                       string
	tr1, tr2                          TimeRange
	expectedResult, expectedRemaining TimeRanges
}

func TestTimeRangeSubtraction(t *testing.T) {
	testCases := []timeRangeSubtractTestCase{
		timeRangeSubtractTestCase{
			description:       "Scenario 1 (s1==s2 && e1==e2)",
			tr1:               mktr("09:00", "10:00"),
			tr2:               mktr("09:00", "10:00"),
			expectedResult:    TimeRanges{},
			expectedRemaining: TimeRanges{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 2 (s1==s2 && e1<e2)",
			tr1:               mktr("09:00", "10:00"),
			tr2:               mktr("09:00", "12:00"),
			expectedResult:    TimeRanges{},
			expectedRemaining: TimeRanges{mktr("10:00", "12:00")},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 3 (s1==s2 && e1>e2)",
			tr1:               mktr("09:00", "12:00"),
			tr2:               mktr("09:00", "10:00"),
			expectedResult:    TimeRanges{mktr("10:00", "12:00")},
			expectedRemaining: TimeRanges{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 4 (s1<s2 && e1==e2)",
			tr1:               mktr("08:00", "12:00"),
			tr2:               mktr("09:00", "12:00"),
			expectedResult:    TimeRanges{mktr("08:00", "09:00")},
			expectedRemaining: TimeRanges{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 5 (s1>s2 && e1==e2)",
			tr1:               mktr("08:00", "12:00"),
			tr2:               mktr("07:00", "12:00"),
			expectedResult:    TimeRanges{},
			expectedRemaining: TimeRanges{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 6 (s1<s2 && e1<e2)",
			tr1:               mktr("07:00", "11:00"),
			tr2:               mktr("08:00", "12:00"),
			expectedResult:    TimeRanges{mktr("07:00", "08:00")},
			expectedRemaining: TimeRanges{mktr("11:00", "12:00")},
		},
		timeRangeSubtractTestCase{
			description: "Scenario 7 (s1<s2 && e1>e2)",
			tr1:         mktr("07:00", "13:00"),
			tr2:         mktr("08:00", "12:00"),
			expectedResult: TimeRanges{
				mktr("07:00", "08:00"),
				mktr("12:00", "13:00"),
			},
			expectedRemaining: TimeRanges{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 8 (s1>s2 && e1<e2)",
			tr1:               mktr("09:00", "11:00"),
			tr2:               mktr("08:00", "12:00"),
			expectedResult:    TimeRanges{},
			expectedRemaining: TimeRanges{mktr("11:00", "12:00")},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 9 (s1>s2 && e1>e2)",
			tr1:               mktr("09:00", "13:00"),
			tr2:               mktr("08:00", "12:00"),
			expectedResult:    TimeRanges{mktr("12:00", "13:00")},
			expectedRemaining: TimeRanges{},
		},
	}

	for _, tc := range testCases {
		diff := tc.tr1.Subtract(tc.tr2)
		assert.Equal(t, diff.result, tc.expectedResult, tc.description)
		assert.Equal(t, diff.remainingSubtractor, tc.expectedRemaining, tc.description)
	}
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
