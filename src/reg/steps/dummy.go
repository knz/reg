package steps

import (
	"reg/t"
)

type stepsource_dummy struct{}

func MakeDummySource() Source {
	return &stepsource_dummy{}
}

func (ts *stepsource_dummy) Start(src <-chan t.Ticks, prod chan<- t.TicksSteps) {
	for ticks := range src {
		prod <- t.TicksSteps{ticks, 0}
	}
}
