// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package t

type Ticks float64
type Steps float64
type Stuff float64
type StuffSteps float64

const (
	SRC_D = 0 // bit 0 unset: DELTAS
	SRC_O = 2 // bit 1 set: SELF ORIGIN
	SRC_M = 3 // bit 0+1 set: MONOTONIC, SELF ORIGIN
	SRC_Z = 4 // bit 2 set: FORCE ORIGIN ZERO
)

type TicksSteps struct {
	Ticks Ticks
	Steps Steps
}

type Status struct {
	Ticks      Ticks
	TicksDelta Ticks
	Steps      Steps
	StepsDelta Steps
	Supply     StuffSteps
	Delta      StuffSteps
}

type Sample struct {
	Ticks Ticks
	Steps Steps
	Usage Stuff
}
