package steps

import (
	"reg/t"
)

func TeeSteps(src <-chan t.TicksSteps, dst chan<- t.TicksSteps, tee chan<- t.Steps) {
	for v := range src {
		dst <- v
		tee <- v.Steps
	}
}
