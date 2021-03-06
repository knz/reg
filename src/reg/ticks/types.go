// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ticks

import "reg/t"

type Source interface {
	Start(src chan<- t.Ticks)
}
