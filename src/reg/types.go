package reg

import (
	"reg/act"
	"reg/sample"
	"reg/steps"
	"reg/t"
	"reg/ticks"
)

type SupplyCmd struct {
	supply t.StuffSteps
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
	Sampler    sample.Sampler
	OutputFile string

	ThrottleType      int
	ThrottleMinPeriod float64

	inputdone chan bool // parse -> wait
}
