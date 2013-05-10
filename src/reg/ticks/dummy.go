package ticks

import (
	"reg/t"
)

type ticksource_dummy struct{}

func (ts *ticksource_dummy) Start(prod chan<- t.Ticks) {
	prod <- t.Ticks(0) // init,
	// then nothing
}

func MakeDummySource() Source {
	return &ticksource_dummy{}
}
