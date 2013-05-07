package ticks

import (
	"reg/t"
)

type ticksource_dummy struct{ ticksource_common }

func (ts *ticksource_dummy) Start() {
	ts.Check()
	go func() {
		ts.source <- t.Ticks(0) // init,
		// then nothing
	}()
}

func MakeDummySource() Source {
	return &ticksource_dummy{}
}
