package reg

import (
	"reg/t"
)

func throttle_ticks(minperiod t.Ticks, src <-chan t.Ticks, prod chan<- t.Ticks) {
	val := float64(0)
	mp := float64(minperiod)

	for v := range src {
		val += float64(v)
		if val >= mp {
			prod <- t.Ticks(val)
			val = 0
		}
	}
}

func mergeticks(ticksext <-chan t.Ticks, ticksctl <-chan t.Ticks, tickssrc chan<- t.Ticks) {

	tickssrc <- <-ticksext // forward init

	for {
		// forward deltas from either external tick source
		// or input stream
		v := t.Ticks(0)
		select {
		case v = <-ticksctl:
		case v = <-ticksext:
		}
		if v > 0 {
			tickssrc <- v
		}
	}
}
