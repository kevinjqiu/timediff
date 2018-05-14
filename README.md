Time Range Subtraction
======================

The code is distributed as a library `github.com/kevinjqiu/timerangesubtraction`

How to Install
--------------

Extract or check out the source code to `$GOPATH/src/github.com/kevinjqiu/timerangesubtraction`

Run Tests
---------

The only dependency for tests is `github.com/stretchr/testify/assert`. Install it using `go get`:

    go get github.com/stretchr/testify/assert

Run tests:

    $ make test
    go test -v
    === RUN   TestTimeRangeSubtraction
    --- PASS: TestTimeRangeSubtraction (0.00s)
    === RUN   TestTimeRangesMerge
    --- PASS: TestTimeRangesMerge (0.00s)
    === RUN   TestTimeRangesSubtraction
    --- PASS: TestTimeRangesSubtraction (0.00s)
    PASS
    ok      github.com/kevinjqiu/timerangesubtraction    0.003s

See coverage report:

    $ make coverage
    go test -coverprofile=cover.out
    PASS
    coverage: 94.0% of statements
    ok      github.com/kevinjqiu/timerangesubtraction       0.003s
    go tool cover -func cover.out
    github.com/kevinjqiu/timerangesubtraction/timerange.go:16:      HasRemainingSubtractor  100.0%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:20:      String                  0.0%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:101:     Subtract                94.7%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:145:     String                  0.0%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:153:     ReplaceAt               100.0%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:168:     Len                     100.0%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:173:     Less                    80.0%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:184:     Swap                    100.0%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:190:     Merge                   100.0%
    github.com/kevinjqiu/timerangesubtraction/timerange.go:214:     Subtract                100.0%
    total:                                                          (statements)            94.0%

Usage
-----

The package exports two structs: `TimeRanges` and `TimeRange`.  `TimeRanges` is a collection of `TimeRange`s.

After define two `TimeRanges` `tr1` and `tr2` you can call `tr1.Subtract(tr2)` to get the time ranges that are in `tr1` but not in `tr2`.

Example:

    const TIME_FORMAT = "15:04"
    trs1 := TimeRanges{
        TimeRange{
            Start: time.Parse(TIME_FORMAT, "09:00"),
            End:   time.Parse(TIME_FORMAT, "11:00"),
        },
        TimeRange{
            Start: time.Parse(TIME_FORMAT, "13:00"),
            End:   time.Parse(TIME_FORMAT, "15:00"),
        },
    }
    trs2 := TimeRanges{
        TimeRange{
            Start: time.Parse(TIME_FORMAT, "09:00"),
            End:   time.Parse(TIME_FORMAT, "09:15"),
        },
        TimeRange{
            Start: time.Parse(TIME_FORMAT, "10:00"),
            End:   time.Parse(TIME_FORMAT, "10:15"),
        },
        TimeRange{
            Start: time.Parse(TIME_FORMAT, "12:30"),
            End:   time.Parse(TIME_FORMAT, "16:00"),
        },
    }

    result := trs1.Subtract(trs2)  // result is a TimeRanges object

How Does It Work
----------------

Let's take the above `TimeRange`s as a demonstration.

    trs1=[(09:00-11:00), (13:00-15:00)] 
    trs2=[(09:00-09:15), (10:00-10:15), (12:30-16:00)]

    trs1.Subtract(trs2)

### Make sure `trs1` and `trs2` are sorted.

We also make the assumption that the `TimeRange`s initially provided do not overlap. If a `TimeRange` `A` is less than `B` if `A.End < B.Start`

Let `i`, `j` be indexes to the time ranges, we enter the loop:

### Iteration 1 (i=0, j=0)

* calculate trs1[i]-trs2[j]

Individual `TimeRange` subtraction is outlined in the comment of `TimeRange.Subtract` method, with illustrations on how the result and remainder are calculated.

In this case: `(09:00-11:00) - (09:00-09:15)`

    trs1=[(09:00-11:00), (13:00-15:00)] 
    trs2=[(09:00-09:15), (10:00-10:15), (12:30-16:00)]
    i=0,j=0
    trs1[i]    09:00-----------------11:00
    trs2[j]    09:00--09:15
    result            09:15----------11:00
    remainder  (empty)

So after the first iteration, the result of the subtraction is `(09:15-11:00)` and there's no remaining `TimeRange` in the subtractor (`trs2[j]`). This means that the subtractor is totally consumed by the subtractee (`trs1[i]`).

* Replace trs1[i] with the result and advance the subtractor index

trs1[0] will be replaced by the result of the previous subtraction, with the index `i` unchanged:

    trs1=[(09:15-11:00), (13:00-15:00)]
                ^
                |
                i

Since the subtractor is completely consumed, we advance the index `j`:

    trs2=[(09:00-09:15), (10:00-10:15), (12:30-16:00)]
                              ^
                              |
                              j

### Iteration 2 (i=0, j=1)

`trs1[0] - trs2[1] = (09:15-11:00) - (10:00-10:15)`

    trs1=[(09:15-11:00), (13:00-15:00)]
    trs2=[(09:00-09:15), (10:00-10:15), (12:30-16:00)]
    i=0,j=1
    trs1[i]    09:15--------------------------11:00
    trs2[j]                10:00---10:15
    result     09:15-------10:00   10:15------11:00
    remainder  (empty)

Note that the result has two spans. trs1[i] will be replaced by both.

    trs1=[(09:15-10:00), (10:15-11:00), (13:00-15:00)]
                ^
                |
                i

Once again, the subtractor is completely consumed so `j` advances:

    trs2=[(09:00-09:15), (10:00-10:15), (12:30-16:00)]
                                              ^
                                              |
                                              j

### Iteration 3 (i=0, j=2)

`trs1[0] - trs2[2] = (09:15-10:00) - (12:30-16:00)`

    i=0,j=1
    trs1[i]    09:15-------10:00
    trs2[j]                       12:30----------16:00
    result     09:15-------10:00
    remainder                     12:30----------16:00

They do not intersect, the result is the same as the subtractee. In this case, we advance the index `i` and keep `j` intact:

    trs1=[(09:15-10:00), (10:15-11:00), (13:00-15:00)]
                               ^
                               |
                               i

### Iteration 4 (i=1, j=2)

`trs1[1] - trs2[2] = (10:15-11:00) - (12:30-16:00)`

    i=0,j=1
    trs1[i]    10:15-------11:00
    trs2[j]                       12:30----------16:00
    result     10:15-------11:00
    remainder                     12:30----------16:00

Again, they do not intersect, the result is the same as the subtractee. again, we advance the index `i` and keep `j` intact:

    trs1=[(09:15-10:00), (10:15-11:00), (13:00-15:00)]
                                              ^
                                              |
                                              i

### Iteration 5 (i=2, j=2)

`trs1[2] - trs2[2] = (13:00-15:00) - (12:30-16:00)`

    i=0,j=1
    trs1[i]             13:00-------15:00
    trs2[j]    12:30--------------------------16:00
    result     (empty)
    remainder  12:30----13:00       15:00-----16:00

Now, the `result` is empty. We replace trs[2] with the empty `TimeRange`, which is equivalent to it being deleted from the slice.

After this iteration, `trs1` is:

    trs1=[(09:15-10:00), (10:15-11:00)]
                                            ^
                                            |
                                            i

### Iteration 6 (i=2, j=2)

`i` now is out of bound for `trs1`, so we break out of the loop.

### Return current trs1

The end result is in trs1. There are cases where subtraction can leave trs1 having overlapping time ranges. To mitigate this, we can `Merge()` before returning so all overlapping or touching time ranges will be merged into a single time range.

In this particular case, we have nothing to merge, so:

    trs1=[(09:15-10:00), (10:15-11:00)]

is the final result.
