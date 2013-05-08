package reg

import (
	"reg/t"
)

func make_status(label string, nres int, ticks *t.Ticks, dticks *t.Ticks,
	steps *t.Steps, dsteps *t.Steps,
	psupply []t.StuffSteps, supply []t.StuffSteps) t.Status {

	v := make([]t.StuffSteps, nres)
	copy(v, supply)
	d := make([]t.StuffSteps, nres)
	for i, x := range v {
		d[i] = x - psupply[i]
	}
	copy(psupply, supply)

	*ticks += *dticks
	*steps += *dsteps

	st := t.Status{
		DomainLabel: label,
		Ticks:       *ticks,
		TicksDelta:  *dticks,
		Steps:       *steps,
		StepsDelta:  *dsteps,
		Supply:      v,
		Delta:       d,
	}
	*dticks = t.Ticks(0)
	*dsteps = t.Steps(0)

	return st
}

func (d *Domain) integrate() {

	dropfirst := true
	nres := len(d.resources)
	supply := make([]t.StuffSteps, nres)

	var qticks, aticks, qdticks, adticks t.Ticks
	var qsteps, asteps, qdsteps, adsteps t.Steps
	qpsupply := make([]t.StuffSteps, nres)
	apsupply := make([]t.StuffSteps, nres)

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
			qdticks += m.ticks
			adticks += m.ticks
			qdsteps += m.steps
			adsteps += m.steps

		case s := <-d.supplycmd:
			supply[s.bin] += s.supply

		case <-d.query:
			d.status <- make_status(d.Label, nres, &qticks, &qdticks, &qsteps, &qdsteps, qpsupply, supply)
		}

		trigger := false
		for i, v := range supply {
			if v < 0 || (v >= 0 && apsupply[i] < 0) {
				trigger = true
			}
		}
		if trigger {
			d.action <- make_status(d.Label, nres, &aticks, &adticks, &asteps, &adsteps, apsupply, supply)
		}

	}
}
