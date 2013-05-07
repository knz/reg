package steps

import (
	"log"
	"reg/t"
)

type stepsource_common struct {
	source chan t.TicksSteps
	ticks  chan t.Ticks
}

func (ts *stepsource_common) SetTicks(ticks chan t.Ticks) {
	ts.ticks = ticks
}
func (ts *stepsource_common) SetSource(src chan t.TicksSteps) {
	ts.source = src
}

func (ts *stepsource_common) Check() {
	if ts.source == nil {
		log.Fatal("no source channel connected")
	}
	if ts.ticks == nil {
		log.Fatal("no ticks channel connected")
	}

}

func TeeSteps(src chan t.TicksSteps, dst chan t.TicksSteps, tee chan t.Steps) {
	go func() {
		for {
			v := <-src
			dst <- v
			tee <- v.Steps
		}
	}()
}
