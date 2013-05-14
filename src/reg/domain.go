// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

	readlines_parse := make(chan string)
	parse_mergeticks := make(chan t.Ticks)
	parse_integrate := make(chan SupplyCmd)
	parse_integrate2 := make(chan bool)
	parse_outmgt := make(chan bool)
	integrate_outmgt := make(chan t.Status)
	integrate_actuator := make(chan t.Status)
	outmgt_integrate := make(chan bool)
	sample_integrate := make(chan t.Sample)
	outmgt_output := make(chan string)
	output_outmgt := make(chan bool)
	tsource_mergeticks := make(chan t.Ticks)
	mergeticks_output := make(chan t.Ticks)
	stepsource_output := make(chan t.TicksSteps)

	go d.TickSource.Start(tsource_mergeticks)
	go mergeticks(tsource_mergeticks, parse_mergeticks, mergeticks_output)

	tickprod := mergeticks_output
	if granularity > 0 {
		throttleticks_output := make(chan t.Ticks)
		go throttle_ticks(granularity, tickprod, throttleticks_output)
		tickprod = throttleticks_output
	}
	if outputpertype == OUTPUT_THROTTLE_TICKS {
		teeticks_throttle := make(chan float64)
		teeticks_output := make(chan t.Ticks)
		go teeticks(tickprod, teeticks_output, teeticks_throttle)
		if ThrottleMinPeriod > 0 {
			go throttle(ThrottleMinPeriod, teeticks_throttle, parse_outmgt)
		} else {
			go forwardctl(teeticks_throttle, parse_outmgt)
		}
		tickprod = teeticks_output
	}

	go d.StepSource.Start(tickprod, stepsource_output)

	stepprod := stepsource_output
	if outputpertype == OUTPUT_THROTTLE_STEPS {
		teesteps_throttle := make(chan float64)
		teesteps_output := make(chan t.TicksSteps)
		go teesteps(stepprod, teesteps_output, teesteps_throttle)
		if ThrottleMinPeriod > 0 {
			go throttle(ThrottleMinPeriod, teesteps_throttle, parse_outmgt)
		} else {
			go forwardctl(teesteps_throttle, parse_outmgt)
		}
		stepprod = teesteps_output
	}
	go d.Sampler.Start(stepprod, sample_integrate)

	go readlines(inputfile, readlines_parse, d.inputdone)
	go parse(readlines_parse, parse_mergeticks, parse_integrate, parse_integrate2, parse_outmgt)

	go d.integrate(integrate_outmgt, integrate_actuator, parse_integrate, parse_integrate2, outmgt_integrate, sample_integrate)

	go d.Actuator.Start(integrate_actuator)

	flood := (outputpertype == OUTPUT_FLOOD)
	go outmgt(flood, parse_outmgt, integrate_outmgt, outmgt_integrate, outmgt_output, output_outmgt)
	go output(outputfile.Fd(), outmgt_output, output_outmgt)

}

func (d *Domain) Wait() {
	<-d.inputdone
}
