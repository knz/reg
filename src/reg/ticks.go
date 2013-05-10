package reg

import "reg/t"

func mergeticks(ticksext <-chan t.Ticks, ticksctl <-chan t.Ticks, tickssrc chan<- t.Ticks) {
	tickssrc <- <-ticksext // forward init
	for {
		// forward deltas from either external tick source
		// or input stream
		val := t.Ticks(0)
		select {
		case val = <-ticksctl:
		case val = <-ticksext:
		}
		if val > 0 {
			tickssrc <- val
		}
	}
}
