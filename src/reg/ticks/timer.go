package ticks

import (
	"reg/t"
	"time"
)

type ticksource_timer struct {
	period time.Duration
}

func MakeTimerSource(period time.Duration) Source {
	return &ticksource_timer{period: period}
}

func (ts *ticksource_timer) Start(prod chan<- t.Ticks) {

	ticker := time.NewTicker(ts.period)
	src := (*ticker).C
	div := 1 / float64(ts.period)

	var lasttime time.Time

	for {
		v := <-src
		d := v.Sub(lasttime)
		lasttime = v
		prod <- t.Ticks(float64(d) * div)
	}
}
