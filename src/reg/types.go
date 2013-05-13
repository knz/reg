// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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
	OUTPUT_EXPLICIT_ONLY = iota
	OUTPUT_FLOOD
	OUTPUT_THROTTLE_STEPS
	OUTPUT_THROTTLE_TICKS
)

type Domain struct {
	TickSource ticks.Source
	StepSource steps.Source
	Actuator   act.Actuator
	Sampler    sample.Sampler

	inputdone chan bool // parse -> wait
}
