// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sample

import "reg/t"

type sampler_dummy struct{ v t.Stuff }

func MakeDummySampler(v t.Stuff) Sampler {
	return &sampler_dummy{v}
}

func (s *sampler_dummy) Start(src <-chan t.TicksSteps, prod chan<- t.Sample) {
	for ts := range src {
		prod <- t.Sample{ts.Ticks, ts.Steps, s.v}
	}
}
