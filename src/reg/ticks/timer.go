package ticks

import ("time"; "log"; "reg/t")

type ticksource_timer struct {
	ticksource_common
	period time.Duration
}

func MakeTimerSource(period time.Duration) Source {
	return &ticksource_timer{period:period}
}

func (ts *ticksource_timer) Start() {
	if ts.source == nil { log.Fatal("no source channel connected") }

	go func() {
		ts.source <- t.Ticks(0) // initial
		for {
			time.Sleep(ts.period)
			ts.source <- t.Ticks(1) // delta
		}
	}()
}
