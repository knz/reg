package steps

import "reg/t"

type Source interface {
	Start()
	SetTicks(src chan t.Ticks)
	SetSource(src chan t.TicksSteps)
}
