// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ticks

import (
	"reg/t"
)

func TeeTicks(dst1 chan t.Ticks, dst2 chan float64, src chan t.Ticks) {
	for a := range src {
		dst1 <- a
		dst2 <- float64(a)
	}
}
