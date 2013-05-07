package reg

func (d *Domain) throttle() {
	val := float64(0)
	for {
		select {
		case ticks := <-d.ticksper:
			if d.ThrottleType != ThrottleTicks {
				continue
			}
			val += float64(ticks)
		case steps := <-d.stepsper:
			if d.ThrottleType != ThrottleSteps {
				continue
			}
			val += float64(steps)
		}
		if val >= d.ThrottleMinPeriod {
			d.statusctl <- true
			val = 0
		}
	}
}
