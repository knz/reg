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

		input:       make(chan string),
		supplycmd:   make(chan SupplyCmd),
		measure:     make(chan Sample),
		query:       make(chan bool),
		status:      make(chan t.Status),
		action:      make(chan t.Status),
		ticksctl:    make(chan t.Ticks),
		statusctl:   make(chan bool),
		tickssrc:    make(chan t.Ticks),
		ticksext:    make(chan t.Ticks),
		ticksin:     make(chan t.Ticks),
		ticksper:    make(chan t.Ticks),
		tickssteps:  make(chan t.TicksSteps),
		tickssteps1: make(chan t.TicksSteps),
		stepsper:    make(chan t.Steps),
		out:         make(chan string),
		outready:    make(chan bool),
		inputdone:   make(chan bool)}

	return &dom
}

func (d *Domain) Start(input io.Reader) {
	d.TickSource.SetSource(d.ticksext)
	d.StepSource.SetTicks(d.tickssrc)
	d.StepSource.SetSource(d.tickssteps1)
	d.Actuator.SetInput(d.action)
	steps.TeeSteps(d.tickssteps1, d.tickssteps, d.stepsper)
	ticks.TeeTicks(d.ticksin, d.ticksper, d.tickssrc)
	d.TickSource.Start()
	d.StepSource.Start()
	d.Actuator.Start()
	go d.readlines(input)
	go d.parse()
	go d.integrate()
	go d.ticksource()
	go d.sample()
	go d.throttle()
	go d.outmgt()
	go d.output()
}

func (d *Domain) Wait() {
	<-d.inputdone
}

func (d *Domain) AddResource(label string, cmd string) {
	resnum := len(d.resources)
	d.resources[resnum] = Resource{label: label, cmd: cmd}
}
