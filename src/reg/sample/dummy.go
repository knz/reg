package sample

import "reg/t"

type sampler_dummy struct{}

func MakeDummySampler() Sampler {
	return &sampler_dummy{}
}

func (s *sampler_dummy) Start(src <-chan t.TicksSteps, prod chan<- t.Sample) {
	for ts := range src {
		prod <- t.Sample{ts.Ticks, ts.Steps, t.Stuff(0)}
	}
}
