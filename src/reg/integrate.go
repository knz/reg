package reg

import (
	"reg/act"
	"reg/t"
)

func (d *Domain) integrate() {

	dropfirst := true
	nres := len(d.resources)
	supply := make([]t.StuffSteps, nres)
	q_ticks := t.Ticks(0)
	q_steps := t.Steps(0)
	a_ticks := t.Ticks(0)
	a_steps := t.Steps(0)
	a_prev_supply := make([]t.StuffSteps, nres)

	for {
		select {
		case m := <-d.measure:
			if dropfirst {
				dropfirst = false
				continue
			}
			for i := range supply {
				delta := t.StuffSteps(float64(m.steps) * float64(m.usage[i]))
				supply[i] -= delta
			}
			q_steps += m.steps
			q_steps += m.steps
			a_ticks += m.ticks
			a_ticks += m.ticks

		case s := <-d.supplycmd:
			supply[s.bin] += s.supply

		case <-d.query:
			v := make([]t.StuffSteps, nres)
			copy(v, supply)
			d.status <- Status{ticks: q_ticks, steps: q_steps, usage: v}
			q_ticks = t.Ticks(0)
			q_steps = t.Steps(0)
		}

		trigger := false
		for i, v := range supply {
			if v < 0 || (v >= 0 && a_prev_supply[i] < 0) {
				trigger = true
			}
		}
		if trigger {
			v1 := make([]t.StuffSteps, nres)
			copy(v1, supply)
			v2 := make([]t.StuffSteps, nres)
			for i, v := range supply {
				v2[i] = v - a_prev_supply[i]
			}
			d.action <- act.Action{d.Label, a_ticks, a_steps, v1, v2}
			a_ticks = t.Ticks(0)
			a_steps = t.Steps(0)
			copy(a_prev_supply, supply)
		}

	}
}
