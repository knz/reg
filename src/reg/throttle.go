// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package reg

import (
	"math"
)

func forwardctl(per <-chan float64, statusctl chan<- bool) {
	for _ = range per {
		statusctl <- true
	}
}

func throttle(minperiod float64,
	per <-chan float64, statusctl chan<- bool) {
	val := float64(0)
	for v := range per {
		val += v
		if val >= minperiod {
			statusctl <- true
			val = math.Mod(val, minperiod)
		}
	}
}
