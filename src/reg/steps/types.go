package steps

import "reg/t"

type Source interface {
	Start(src <-chan t.Ticks, prod chan<- t.TicksSteps)
}
