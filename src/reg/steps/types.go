package steps

import "reg/t"

type TicksSteps struct {
	ticks t.Ticks
	steps t.Steps
}

type Source interface {
	Start()
	SetTicks(src chan t.Ticks)
	SetSource(src chan TicksSteps)
}
