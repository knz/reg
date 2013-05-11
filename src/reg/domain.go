package reg

import (
	"io"
	"reg/act"
	"reg/sample"
	"reg/steps"
	"reg/t"
	"reg/ticks"
)

func MakeDomain(label string, ts ticks.Source, ss steps.Source, actuator act.Actuator,
	sampler sample.Sampler) *Domain {
	dom := Domain{
		Label:      label,
		TickSource: ts,
		StepSource: ss,
		Actuator:   actuator,
		Sampler:    sampler,

		inputdone: make(chan bool)}

	return &dom
}

func (d *Domain) Start(input io.Reader) {

	tsource_mergeticks := make(chan t.Ticks)
	go d.TickSource.Start(tsource_mergeticks)

	readlines_parse := make(chan string)
	go readlines(input, readlines_parse, d.inputdone)

	parse_mergeticks := make(chan t.Ticks)
	parse_integrate := make(chan SupplyCmd)
	parse_outmgt := make(chan bool)
	go parse(readlines_parse, parse_mergeticks, parse_integrate, parse_outmgt)

	integrate_outmgt := make(chan t.Status)
	integrate_actuator := make(chan t.Status)
	outmgt_integrate := make(chan bool)
	sample_integrate := make(chan t.Sample)
	go d.integrate(integrate_outmgt, integrate_actuator, parse_integrate, outmgt_integrate, sample_integrate)

	go d.Actuator.Start(integrate_actuator)

	mergeticks_teeticks := make(chan t.Ticks)
	go mergeticks(tsource_mergeticks, parse_mergeticks, mergeticks_teeticks)

	teeticks_ssource := make(chan t.Ticks)
	teeticks_throttle := make(chan t.Ticks)
	go ticks.TeeTicks(teeticks_ssource, teeticks_throttle, mergeticks_teeticks)

	ssource_teesteps := make(chan t.TicksSteps)
	go d.StepSource.Start(teeticks_ssource, ssource_teesteps)

	teesteps_sample := make(chan t.TicksSteps)
	teesteps_throttle := make(chan t.Steps)
	go steps.TeeSteps(ssource_teesteps, teesteps_sample, teesteps_throttle)

	go d.Sampler.Start(teesteps_sample, sample_integrate)

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
