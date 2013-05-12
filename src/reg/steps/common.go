package steps

import (
	"reg/t"
)

func TeeSteps(src <-chan t.TicksSteps, dst chan<- t.TicksSteps, tee chan<- float64) {
	for v := range src {
		dst <- v
		tee <- float64(v.Steps)
	}
}
