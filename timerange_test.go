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
	description       string
	tr1, tr2          TimeRange
	expectedResult    TimeRanges
	expectedRemaining TimeRange
}

func TestTimeRangeSubtraction(t *testing.T) {
	testCases := []timeRangeSubtractTestCase{
		timeRangeSubtractTestCase{
			description:       "Scenario 1 (s1==s2 && e1==e2)",
			tr1:               mktr("09:00", "10:00"),
			tr2:               mktr("09:00", "10:00"),
			expectedResult:    TimeRanges{},
			expectedRemaining: TimeRange{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 2 (s1==s2 && e1<e2)",
			tr1:               mktr("09:00", "10:00"),
			tr2:               mktr("09:00", "12:00"),
			expectedResult:    TimeRanges{},
			expectedRemaining: mktr("10:00", "12:00"),
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 3 (s1==s2 && e1>e2)",
			tr1:               mktr("09:00", "12:00"),
			tr2:               mktr("09:00", "10:00"),
			expectedResult:    TimeRanges{mktr("10:00", "12:00")},
			expectedRemaining: TimeRange{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 4 (s1<s2 && e1==e2)",
			tr1:               mktr("08:00", "12:00"),
			tr2:               mktr("09:00", "12:00"),
			expectedResult:    TimeRanges{mktr("08:00", "09:00")},
			expectedRemaining: TimeRange{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 5 (s1>s2 && e1==e2)",
			tr1:               mktr("08:00", "12:00"),
			tr2:               mktr("07:00", "12:00"),
			expectedResult:    TimeRanges{},
			expectedRemaining: TimeRange{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 6 (s1<s2 && e1<e2)",
			tr1:               mktr("07:00", "11:00"),
			tr2:               mktr("08:00", "12:00"),
			expectedResult:    TimeRanges{mktr("07:00", "08:00")},
			expectedRemaining: mktr("11:00", "12:00"),
		},
		timeRangeSubtractTestCase{
			description: "Scenario 7 (s1<s2 && e1>e2)",
			tr1:         mktr("07:00", "13:00"),
			tr2:         mktr("08:00", "12:00"),
			expectedResult: TimeRanges{
				mktr("07:00", "08:00"),
				mktr("12:00", "13:00"),
			},
			expectedRemaining: TimeRange{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 8 (s1>s2 && e1<e2)",
			tr1:               mktr("09:00", "11:00"),
			tr2:               mktr("08:00", "12:00"),
			expectedResult:    TimeRanges{},
			expectedRemaining: mktr("11:00", "12:00"),
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 9 (s1>s2 && e1>e2)",
			tr1:               mktr("09:00", "13:00"),
			tr2:               mktr("08:00", "12:00"),
			expectedResult:    TimeRanges{mktr("12:00", "13:00")},
			expectedRemaining: TimeRange{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 10 (s1>=e2)",
			tr1:               mktr("09:00", "13:00"),
			tr2:               mktr("05:00", "07:00"),
			expectedResult:    TimeRanges{mktr("09:00", "13:00")},
			expectedRemaining: TimeRange{},
		},
		timeRangeSubtractTestCase{
			description:       "Scenario 11 (s2>=e1)",
			tr1:               mktr("09:00", "13:00"),
			tr2:               mktr("15:00", "17:00"),
			expectedResult:    TimeRanges{mktr("09:00", "13:00")},
			expectedRemaining: mktr("15:00", "17:00"),
		},
	}

	for _, tc := range testCases {
		diff := tc.tr1.Subtract(tc.tr2)
		assert.Equal(t, diff.result, tc.expectedResult, tc.description)
		assert.Equal(t, diff.remainingSubtractor, tc.expectedRemaining, tc.description)
	}
}

type timeRangeMergeTestCase struct {
	description    string
	timeRanges     TimeRanges
	expectedResult TimeRanges
}

func TestTimeRangesMerge(t *testing.T) {
	testCases := []timeRangeMergeTestCase{
		timeRangeMergeTestCase{
			description:    "no time ranges",
			timeRanges:     TimeRanges{},
			expectedResult: TimeRanges{},
		},
		timeRangeMergeTestCase{
			description: "only one time range",
			timeRanges: TimeRanges{
				mktr("09:00", "10:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:00", "10:00"),
			},
		},
		timeRangeMergeTestCase{
			description: "no overlapping time ranges",
			timeRanges: TimeRanges{
				mktr("09:00", "10:00"),
				mktr("11:00", "12:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:00", "10:00"),
				mktr("11:00", "12:00"),
			},
		},
		timeRangeMergeTestCase{
			description: "connected time ranges",
			timeRanges: TimeRanges{
				mktr("09:00", "11:00"),
				mktr("11:00", "12:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:00", "12:00"),
			},
		},
		timeRangeMergeTestCase{
			description: "equal time ranges",
			timeRanges: TimeRanges{
				mktr("09:00", "11:00"),
				mktr("09:00", "11:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:00", "11:00"),
			},
		},
		timeRangeMergeTestCase{
			description: "second time range falls completely within the first one",
			timeRanges: TimeRanges{
				mktr("09:00", "12:00"),
				mktr("09:30", "11:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:00", "12:00"),
			},
		},
		timeRangeMergeTestCase{
			description: "overlapping time ranges",
			timeRanges: TimeRanges{
				mktr("09:00", "12:00"),
				mktr("09:45", "13:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:00", "13:00"),
			},
		},
		timeRangeMergeTestCase{
			description: "multiple time ranges",
			timeRanges: TimeRanges{
				mktr("09:00", "12:00"),
				mktr("09:45", "13:00"),
				mktr("10:45", "15:00"),
				mktr("17:45", "18:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:00", "15:00"),
				mktr("17:45", "18:00"),
			},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedResult, tc.timeRanges.Merge(), tc.description)
	}
}

type TimeRangesSubtractTestCase struct {
	description    string
	tr1, tr2       TimeRanges
	expectedResult TimeRanges
}

func TestTimeRangesSubtraction(t *testing.T) {
	testCases := []TimeRangesSubtractTestCase{
		TimeRangesSubtractTestCase{
			description:    "tr1 supersedes tr2",
			tr1:            TimeRanges{mktr("09:00", "10:00")},
			tr2:            TimeRanges{mktr("09:00", "09:30")},
			expectedResult: TimeRanges{mktr("09:30", "10:00")},
		},
		TimeRangesSubtractTestCase{
			description:    "tr1 equals tr2",
			tr1:            TimeRanges{mktr("09:30", "10:30")},
			tr2:            TimeRanges{mktr("09:30", "10:30")},
			expectedResult: TimeRanges{},
		},
		TimeRangesSubtractTestCase{
			description:    "tr1 borders tr2",
			tr1:            TimeRanges{mktr("09:00", "09:30")},
			tr2:            TimeRanges{mktr("09:30", "15:00")},
			expectedResult: TimeRanges{mktr("09:00", "09:30")},
		},
		TimeRangesSubtractTestCase{
			description:    "tr1 and tr2 do not intersect",
			tr1:            TimeRanges{mktr("09:00", "09:30")},
			tr2:            TimeRanges{mktr("09:31", "15:00")},
			expectedResult: TimeRanges{mktr("09:00", "09:30")},
		},
		TimeRangesSubtractTestCase{
			description: "multiple overlapping time ranges",
			tr1: TimeRanges{
				mktr("09:00", "09:30"),
				mktr("10:00", "10:30"),
			},
			tr2: TimeRanges{
				mktr("09:15", "10:15"),
			},
			expectedResult: TimeRanges{
				mktr("09:00", "09:15"),
				mktr("10:15", "10:30"),
			},
		},
		TimeRangesSubtractTestCase{
			description: "multiple overlapping time ranges",
			tr1: TimeRanges{
				mktr("09:00", "11:00"),
				mktr("13:00", "15:00"),
			},
			tr2: TimeRanges{
				mktr("09:00", "09:15"),
				mktr("10:00", "10:15"),
				mktr("12:30", "16:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:15", "10:00"),
				mktr("10:15", "11:00"),
			},
		},
		TimeRangesSubtractTestCase{
			description: "multiple time ranges subtraction",
			tr1: TimeRanges{
				mktr("09:00", "11:00"),
				mktr("13:00", "15:00"),
			},
			tr2: TimeRanges{
				mktr("09:00", "09:15"),
				mktr("10:00", "10:15"),
				mktr("12:30", "16:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:15", "10:00"),
				mktr("10:15", "11:00"),
			},
		},
		TimeRangesSubtractTestCase{
			description: "multiple out-of-order time ranges subtraction",
			tr1: TimeRanges{
				mktr("13:00", "15:00"),
				mktr("09:00", "11:00"),
			},
			tr2: TimeRanges{
				mktr("10:00", "10:15"),
				mktr("09:00", "09:15"),
				mktr("12:30", "16:00"),
			},
			expectedResult: TimeRanges{
				mktr("09:15", "10:00"),
				mktr("10:15", "11:00"),
			},
		},
	}

	for _, tc := range testCases {
		trdiff := tc.tr1.Subtract(tc.tr2)
		assert.Equal(t, tc.expectedResult, trdiff, tc.description)
	}
}
