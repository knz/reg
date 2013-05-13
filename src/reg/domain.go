package reg

import (
	"os"
	"reg/act"
	"reg/sample"
	"reg/steps"
	"reg/t"
	"reg/ticks"
)

func MakeDomain(ts ticks.Source,
	ss steps.Source,
	actuator act.Actuator,
	sampler sample.Sampler) *Domain {
	dom := Domain{
		TickSource: ts,
		StepSource: ss,
		Actuator:   actuator,
		Sampler:    sampler,

		inputdone: make(chan bool)}

	return &dom
}

func (d *Domain) Start(inputfile *os.File, outputfile *os.File,
	outputpertype int, ThrottleMinPeriod float64, granularity t.Ticks) {

	tsource_mergeticks := make(chan t.Ticks)
	readlines_parse := make(chan string)
	parse_mergeticks := make(chan t.Ticks)
	parse_integrate := make(chan SupplyCmd)
	parse_outmgt := make(chan bool)
	integrate_outmgt := make(chan t.Status)
	integrate_actuator := make(chan t.Status)
	outmgt_integrate := make(chan bool)
	sample_integrate := make(chan t.Sample)
	outmgt_output := make(chan string)
	output_outmgt := make(chan bool)

	mergeticksOutput := make(chan t.Ticks)
	mergeticksOutputT := mergeticksOutput

	if granularity > 0 {
		mergeticksOutputT = make(chan t.Ticks)
		go throttle_ticks(granularity, mergeticksOutput, mergeticksOutputT)
	}

	ssourceInput := mergeticksOutputT
	ssourceOutput := make(chan t.TicksSteps)
	sampleInput := ssourceOutput
	flood := false

	switch outputpertype {
	case OUTPUT_FLOOD:
		flood = true
	case OUTPUT_THROTTLE_STEPS:
		teesteps_throttle := make(chan float64)
		sampleInput = make(chan t.TicksSteps)
		go steps.TeeSteps(ssourceOutput, sampleInput, teesteps_throttle)
		if ThrottleMinPeriod > 0 {
			go throttle(ThrottleMinPeriod, teesteps_throttle, parse_outmgt)
		} else {
			go forwardctl(teesteps_throttle, parse_outmgt)
		}
	case OUTPUT_THROTTLE_TICKS:
		teeticks_throttle := make(chan float64)
		ssourceInput = make(chan t.Ticks)
		go ticks.TeeTicks(ssourceInput, teeticks_throttle, mergeticksOutputT)
		if ThrottleMinPeriod > 0 {
			go throttle(ThrottleMinPeriod, teeticks_throttle, parse_outmgt)
		} else {
			go forwardctl(teeticks_throttle, parse_outmgt)
		}
	}

	go readlines(inputfile, readlines_parse, d.inputdone)
	go parse(readlines_parse, parse_mergeticks, parse_integrate, parse_outmgt)

	go d.TickSource.Start(tsource_mergeticks)
	go mergeticks(tsource_mergeticks, parse_mergeticks, mergeticksOutput)
	go d.StepSource.Start(ssourceInput, ssourceOutput)
	go d.Sampler.Start(sampleInput, sample_integrate)

	go d.integrate(integrate_outmgt, integrate_actuator, parse_integrate, outmgt_integrate, sample_integrate)

	go d.Actuator.Start(integrate_actuator)

	go outmgt(flood, parse_outmgt, integrate_outmgt, outmgt_integrate, outmgt_output, output_outmgt)
	go output(outputfile.Fd(), outmgt_output, output_outmgt)

}

func (d *Domain) Wait() {
	<-d.inputdone
}
