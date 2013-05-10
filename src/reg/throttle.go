package reg

import (
	"math"
	"reg/t"
)

func (d *Domain) throttle(ticksper <-chan t.Ticks, stepsper <-chan t.Steps, statusctl chan<- bool) {
	val := float64(0)
	for {
		select {
		case ticks := <-ticksper:
			if d.ThrottleType != ThrottleTicks {
				continue
			}
			val += float64(ticks)
		case steps := <-stepsper:
			if d.ThrottleType != ThrottleSteps {
				continue
			}
			val += float64(steps)
		}
		if val >= d.ThrottleMinPeriod {
			statusctl <- true
			val = math.Mod(val, d.ThrottleMinPeriod)
		}
	}
}
