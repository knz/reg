package reg

import (
	"reg/t"
)

type qstate struct {
	ticks   t.Ticks
	dticks  t.Ticks
	steps   t.Steps
	dsteps  t.Steps
	psupply t.StuffSteps
}

func (qs *qstate) make_status(supply t.StuffSteps) t.Status {

	qs.ticks += qs.dticks
	qs.steps += qs.dsteps

	st := t.Status{
		Ticks:      qs.ticks,
		TicksDelta: qs.dticks,
		Steps:      qs.steps,
		StepsDelta: qs.dsteps,
		Supply:     supply,
		Delta:      supply - qs.psupply,
	}
	qs.psupply = supply
	qs.dticks = t.Ticks(0)
	qs.dsteps = t.Steps(0)

	return st
}

func (d *Domain) integrate(
	dropfirst bool,
	status chan<- t.Status,
	action chan<- t.Status,
	supplycmd <-chan SupplyCmd,
	query <-chan bool,
	measure <-chan t.Sample) {

	var supply t.StuffSteps
	var as, qs qstate

	for {
		update := false
		select {
		case m := <-measure:
			if dropfirst {
				dropfirst = false
				continue
			}

			delta := t.StuffSteps(float64(m.Steps) * float64(m.Usage))
			supply -= delta
			update = true

			qs.dticks += m.Ticks
			as.dticks += m.Ticks
			qs.dsteps += m.Steps
			as.dsteps += m.Steps

		case s := <-supplycmd:
			supply += s.supply
			update = true

		case <-query:
			status <- qs.make_status(supply)
		}

		if update && (supply < 0 || (supply >= 0 && as.psupply < 0)) {
			action <- as.make_status(supply)
		}

	}
}
