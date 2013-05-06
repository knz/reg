package reg
import "reg/ticks"
func dupticks(out1 chan<- ticks.Ticks, out2 chan<- ticks.Ticks, in <-chan ticks.Ticks) {
	for a := range in { out1 <- a; out2 <- a; }
}
