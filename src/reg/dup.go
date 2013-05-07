package reg

import "reg/t"

func dupticks(out1 chan<- t.Ticks, out2 chan<- t.Ticks, in <-chan t.Ticks) {
	for a := range in {
		out1 <- a
		out2 <- a
	}
}
