package ticks

import (
	"reg/t"
)

type ticksource_dummy struct{ v t.Ticks }

func (ts *ticksource_dummy) Start(prod chan<- t.Ticks) {
	prod <- t.Ticks(ts.v) // init,
	// then nothing
}

func MakeDummySource(v t.Ticks) Source {
	return &ticksource_dummy{v}
}
