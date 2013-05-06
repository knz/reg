package reg

func dupticks(out1 chan<- Ticks, out2 chan<- Ticks, in <-chan Ticks) {
	for a := range in { out1 <- a; out2 <- a; }
}
