package reg

import (
	"io"
	"reg/act"
	"reg/steps"
	"reg/t"
	"reg/ticks"
)

func MakeDomain(label string, ts ticks.Source, ss steps.Source, actuator act.Actuator) *Domain {
	dom := Domain{
		Label:      label,
		TickSource: ts,
		StepSource: ss,
		Actuator:   actuator,

		resources: make(map[int]Resource),

		inputdone: make(chan bool)}

	return &dom
}

func (d *Domain) Start(input io.Reader) {

	tsource_mergeticks := make(chan t.Ticks)
	d.TickSource.SetSource(tsource_mergeticks)
	d.TickSource.Start()

	readlines_parse := make(chan string)
	go readlines(input, readlines_parse, d.inputdone)

	parse_mergeticks := make(chan t.Ticks)
	parse_integrate := make(chan SupplyCmd)
	parse_outmgt := make(chan bool)
	go parse(readlines_parse, parse_mergeticks, parse_integrate, parse_outmgt)

	integrate_outmgt := make(chan t.Status)
	integrate_actuator := make(chan t.Status)
	outmgt_integrate := make(chan bool)
	sample_integrate := make(chan Sample)
	go d.integrate(integrate_outmgt, integrate_actuator, parse_integrate, outmgt_integrate, sample_integrate)

	d.Actuator.SetInput(integrate_actuator)
	d.Actuator.Start()

	mergeticks_teeticks := make(chan t.Ticks)
	go d.mergeticks(tsource_mergeticks, parse_mergeticks, mergeticks_teeticks)

	teeticks_ssource := make(chan t.Ticks)
	teeticks_throttle := make(chan t.Ticks)
	ticks.TeeTicks(teeticks_ssource, teeticks_throttle, mergeticks_teeticks)

	d.StepSource.SetTicks(teeticks_ssource)
	ssource_teesteps := make(chan t.TicksSteps)
	d.StepSource.SetSource(ssource_teesteps)
	d.StepSource.Start()

	teesteps_sample := make(chan t.TicksSteps)
	teesteps_throttle := make(chan t.Steps)
	steps.TeeSteps(ssource_teesteps, teesteps_sample, teesteps_throttle)

	go d.sample(teesteps_sample, sample_integrate)

	//throttle_outmgt := make(chan bool)
	go d.throttle(teeticks_throttle, teesteps_throttle, parse_outmgt)

	outmgt_output := make(chan string)
	output_outmgt := make(chan bool)
	go d.outmgt(parse_outmgt, integrate_outmgt, outmgt_integrate, outmgt_output, output_outmgt)

	go d.output(outmgt_output, output_outmgt)
}

func (d *Domain) Wait() {
	<-d.inputdone
}

func (d *Domain) AddResource(label string, cmd string) {
	resnum := len(d.resources)
	d.resources[resnum] = Resource{label: label, cmd: cmd}
}
