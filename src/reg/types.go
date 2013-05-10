package reg

import (
	"reg/act"
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

type Resource struct {
	label string
	cmd   string
}

const (
	ThrottleSteps = iota
	ThrottleTicks
)

type Domain struct {
	Label      string
	TickSource ticks.Source
	StepSource steps.Source
	Actuator   act.Actuator
	OutputFile string

	ThrottleType      int
	ThrottleMinPeriod float64

	// Resource management
	resources map[int]Resource

	inputdone chan bool // parse -> wait
}
