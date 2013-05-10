package reg

import (
	"reg/t"
)

func make_status(label string, ticks *t.Ticks, dticks *t.Ticks,
	steps *t.Steps, dsteps *t.Steps,
	psupply *t.StuffSteps, supply *t.StuffSteps) t.Status {

	*ticks += *dticks
	*steps += *dsteps

	st := t.Status{
		DomainLabel: label,
		Ticks:       *ticks,
		TicksDelta:  *dticks,
		Steps:       *steps,
		StepsDelta:  *dsteps,
		Supply:      *supply,
		Delta:       *supply - *psupply,
	}
	*psupply = *supply
	*dticks = t.Ticks(0)
	*dsteps = t.Steps(0)

	return st
}

func (d *Domain) integrate(
	status chan<- t.Status,
	action chan<- t.Status,
	supplycmd <-chan SupplyCmd,
	query <-chan bool,
	measure <-chan t.Sample) {

	dropfirst := true
	var supply t.StuffSteps

	var qticks, aticks, qdticks, adticks t.Ticks
	var qsteps, asteps, qdsteps, adsteps t.Steps
	var qpsupply, apsupply t.StuffSteps

	for {
		select {
		case m := <-measure:
			if dropfirst {
				dropfirst = false
				continue
			}

			delta := t.StuffSteps(float64(m.Steps) * float64(m.Usage))
			supply -= delta

			qdticks += m.Ticks
			adticks += m.Ticks
			qdsteps += m.Steps
			adsteps += m.Steps

		case s := <-supplycmd:
			supply += s.supply

		case <-query:
			status <- make_status(d.Label, &qticks, &qdticks, &qsteps, &qdsteps, &qpsupply, &supply)
		}

		if supply < 0 || (supply >= 0 && apsupply < 0) {
			action <- make_status(d.Label, &aticks, &adticks, &asteps, &adsteps, &apsupply, &supply)
		}

	}
}
