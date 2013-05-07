package reg

import (
	"reg/t"
	"reg/ticks"
)

type TicksSteps struct {
	ticks t.Ticks
	steps t.Steps
}

type SupplyCmd struct {
	bin    int
	supply t.StuffSteps
}

type Sample struct {
	ticks t.Ticks
	steps t.Steps
	usage []t.Stuff
}

type Status struct {
	ticks t.Ticks
	steps t.Steps
	usage []t.StuffSteps
}

type Action struct {
	bin           int
	currentsupply t.StuffSteps
	delta         t.StuffSteps
}

type Resource struct {
	label string
	cmd   string
}

const (
	ThrottleSteps = iota
	ThrottleTicks
)

type Domain struct {
	Label       string
	TickSource  ticks.Source
	StepsCmd    string
	ProtocolCmd string
	OutputFile  string

	ThrottleType      int
	ThrottleMinPeriod float64

	// Resource management
	resources map[int]Resource

	// Channels
	input      chan string     // readlines -> parse
	supplycmd  chan SupplyCmd  // parse -> integrate
	measure    chan Sample     // sample -> integrate
	query      chan bool       // outputmgt -> integrate
	status     chan Status     // integrate -> outputmgt
	action     chan Action     // integrate -> protocol
	ticksctl   chan t.Ticks    // parse -> ticksource
	statusctl  chan bool       // parse -> outmgt
	ticksext   chan t.Ticks    // tickext -> ticksource
	tickssrc   chan t.Ticks    // ticksource -> dup
	ticksin    chan t.Ticks    // dup -> stepsource
	ticksper   chan t.Ticks    // dup -> throttle
	tickssteps chan TicksSteps // stepsource -> sample
	stepsper   chan t.Steps    // stepsource -> throttle

	out      chan string // outmgt -> output
	outready chan bool   // output -> outmgt

	inputdone chan bool // parse -> wait
}
