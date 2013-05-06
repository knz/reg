package reg
import "reg/ticks"
func (d *Domain) ticksource() {
	d.tickssrc <- <- d.ticksext // forward init
	for {
		// forward deltas from either external tick source
		// or input stream
		t := ticks.Ticks(0)
		select {
		case t = <- d.ticksctl:
		case t = <- d.ticksext:
		}
		if t > 0 {
			d.tickssrc <- t
		}
	}
}
