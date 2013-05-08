package act

import "reg/t"

type Action struct {
	DomainLabel string
	TicksDelta  t.Ticks
	StepsDelta  t.Steps
	Supply      []t.StuffSteps
	Delta       []t.StuffSteps
}

type Actuator interface {
	Start()
	SetInput(src chan Action)
}
