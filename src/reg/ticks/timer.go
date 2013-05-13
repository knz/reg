package ticks

import (
	"reg/t"
	"time"
)

type ticksource_timer struct {
	period time.Duration
	stype  int
	divper bool
}

func MakeTimerSource(period time.Duration, stype int, divper bool) Source {
	return &ticksource_timer{period, stype, divper}
}

func (ts *ticksource_timer) Start(prod chan<- t.Ticks) {

	ticker := time.NewTicker(ts.period)
	src := (*ticker).C

	var div float64

	if ts.divper {
		div = 1 / float64(ts.period)
	} else {
		div = 1 / 1e9
	}

	var lasttime time.Time

	if ts.stype&t.SRC_Z != 0 {
		// origin forced to zero, so consume
		// first tick to define origin
		lasttime = <-src
	}

	for v := range src {
		d := v.Sub(lasttime)
		lasttime = v
		prod <- t.Ticks(float64(d) * div)
	}
}
