package steps

import (
	"reg/t"
)

type stepsource_dummy struct{ v t.Steps }

func MakeDummySource(v t.Steps) Source {
	return &stepsource_dummy{v}
}

func (ts *stepsource_dummy) Start(src <-chan t.Ticks, prod chan<- t.TicksSteps) {
	for ticks := range src {
		prod <- t.TicksSteps{ticks, ts.v}
	}
}
