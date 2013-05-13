// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package steps

import (
	"reg/t"
)

func TeeSteps(src <-chan t.TicksSteps, dst chan<- t.TicksSteps, tee chan<- float64) {
	for v := range src {
		dst <- v
		tee <- float64(v.Steps)
	}
}
