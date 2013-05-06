package reg
import "reg/t"
func (d *Domain) integrate() {

	steps := t.Steps(0)
	ticks := t.Ticks(0)
	nres := len(d.resources)
	supply := make([]t.StuffSteps, nres)
	deltas := make([]t.StuffSteps, nres)
	trigger := make([]bool, nres)

	for {
		for i := range trigger { trigger[i] = false }
		select {
		case m := <- d.measure:
			for i := range supply {
				deltas[i] = t.StuffSteps(float64(m.steps) * float64(m.usage[i]))
				v := supply[i] - deltas[i]
				if (v > 0 && supply[i] <= 0) || v <= 0 { trigger[i] = true }
				supply[i] = v
			}
			steps += m.steps
			ticks += m.ticks

		case s := <- d.supplycmd:
			v := supply[s.bin] + s.supply
			if (v > 0 && supply[s.bin] <= 0) || v <= 0 { trigger[s.bin] = true }
			supply[s.bin] = v

		case <- d.query:
			v := make([]t.StuffSteps, nres)
			copy(v, supply)
			d.status <- Status{ ticks : ticks, steps : steps, usage : v }
		}

		for i, t := range trigger {
			if t {
				d.action <- Action { bin : i, currentsupply : supply[i], delta : deltas[i] }
			}
		}
	}
}