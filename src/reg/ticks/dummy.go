// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ticks

import (
	"reg/t"
)

type ticksource_dummy struct{ v t.Ticks }

func (ts *ticksource_dummy) Start(prod chan<- t.Ticks) {
	prod <- t.Ticks(ts.v) // init,
	// then nothing
}

func MakeDummySource(v t.Ticks) Source {
	return &ticksource_dummy{v}
}
