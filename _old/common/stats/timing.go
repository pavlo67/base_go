package stats

import (
	"fmt"
	"time"
)

type Timing struct {
	Cnt int
	Sum time.Duration
	Min time.Duration
	Max time.Duration
}

func (tm *Timing) Add(momentPrev time.Time) {
	if tm == nil {
		return
	}

	d := time.Now().Sub(momentPrev)

	if tm.Cnt == 0 {
		tm.Cnt = 1
		tm.Sum = d
		tm.Min = d
		tm.Max = d
	} else {
		tm.Cnt++
		tm.Sum += d
		if d < tm.Min {
			tm.Min = d
		}
		if d > tm.Max {
			tm.Max = d
		}
	}
}

func (tm *Timing) String() string {
	if tm == nil || tm.Cnt <= 0 {
		return "<empty tm>"
	}

	return fmt.Sprintf("tm: cnt = %d, avg = %s, min = %s, max = %s",
		tm.Cnt, tm.Sum/time.Duration(tm.Cnt), tm.Min, tm.Max)
}
