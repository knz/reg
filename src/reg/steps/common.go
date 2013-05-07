package steps

import "reg/t"

type stepsource_common struct {
	source chan TicksSteps
	ticks  chan t.Ticks
}

func (ts *stepsource_common) SetTicks(ticks chan t.Ticks) {
	ts.ticks = ticks
}
func (ts *stepsource_common) SetSource(src chan t.Steps) {
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

func TeeSteps(src chan TicksSteps, dst chan TicksSteps, tee chan Steps) {
	go func() {
		for {
			v := <-src
			dst <- v
			tee <- v.steps
		}
	}()
}
