package reg

import "reg/t"

func (d *Domain) ticksource() {
	d.tickssrc <- <- d.ticksext // forward init
	for {
		// forward deltas from either external tick source
		// or input stream
		val := t.Ticks(0)
		select {
		case val = <- d.ticksctl:
		case val = <- d.ticksext:
		}
		if val > 0 {
			d.tickssrc <- val
		}
	}
}
