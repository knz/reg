// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sample

import "reg/t"

type Sampler interface {
	Start(src <-chan t.TicksSteps, prod chan<- t.Sample)
}
