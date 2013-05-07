package reg

import (
	"reg/steps"
	"reg/t"
	"reg/ticks"
)

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
	StepSource  steps.Source
	ProtocolCmd string
	OutputFile  string

	ThrottleType      int
	ThrottleMinPeriod float64

	// Resource management
	resources map[int]Resource

	// Channels
	input       chan string       // readlines -> parse
	supplycmd   chan SupplyCmd    // parse -> integrate
	measure     chan Sample       // sample -> integrate
	query       chan bool         // outputmgt -> integrate
	status      chan Status       // integrate -> outputmgt
	action      chan Action       // integrate -> protocol
	ticksctl    chan t.Ticks      // parse -> ticksource
	statusctl   chan bool         // parse -> outmgt
	ticksext    chan t.Ticks      // tickext -> ticksource
	tickssrc    chan t.Ticks      // ticksource -> dup
	ticksin     chan t.Ticks      // dup -> stepsource
	ticksper    chan t.Ticks      // dup -> throttle
	tickssteps1 chan t.TicksSteps // stepsource -> teesteps
	tickssteps  chan t.TicksSteps // teesteps -> sample
	stepsper    chan t.Steps      // teesteps -> throttle

	out      chan string // outmgt -> output
	outready chan bool   // output -> outmgt

	inputdone chan bool // parse -> wait
}
