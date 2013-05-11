package sample

import "reg/t"

type Sampler interface {
	Start(src <-chan t.TicksSteps, prod chan<- t.Sample)
}
