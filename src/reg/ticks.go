// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package reg

import (
	"reg/t"
)

func throttle_ticks(minperiod t.Ticks, src <-chan t.Ticks, prod chan<- t.Ticks) {
	val := float64(0)
	mp := float64(minperiod)

	for v := range src {
		val += float64(v)
		if val >= mp {
			prod <- t.Ticks(val)
			val = 0
		}
	}
}

func teeticks(src <-chan t.Ticks, dst1 chan<- t.Ticks, dst2 chan<- float64) {
	for a := range src {
		dst1 <- a
		dst2 <- float64(a)
	}
}

func teesteps(src <-chan t.TicksSteps, dst chan<- t.TicksSteps, tee chan<- float64) {
	for v := range src {
		dst <- v
		tee <- float64(v.Steps)
	}
}

func mergeticks(src_ext <-chan t.Ticks, src_ctl <-chan t.Ticks, prod chan<- t.Ticks) {

	prod <- <-src_ext // forward init

	for {
		// forward deltas from either external tick source
		// or input stream
		v := t.Ticks(0)
		select {
		case v = <-src_ctl:
		case v = <-src_ext:
		}
		if v > 0 {
			prod <- v
		}
	}
}
