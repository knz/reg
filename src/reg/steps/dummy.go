package steps

import (
	"reg/t"
)

type stepsource_dummy struct{ stepsource_common }

func (ts *stepsource_dummy) Start() {
	ts.Check()

	go func() {
		for {
			ts.source <- TicksSteps{<-ts.ticks, 0}
		}
	}()
}
