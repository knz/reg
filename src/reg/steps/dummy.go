package steps

import (
	"reg/t"
)

type stepsource_dummy struct{ stepsource_common }

func MakeDummySource() Source {
	return &stepsource_dummy{}
}

func (ts *stepsource_dummy) Start() {
	ts.Check()

	go func() {
		for {
			ts.source <- t.TicksSteps{<-ts.ticks, 0}
		}
	}()
}
