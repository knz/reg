package ticks

import "reg/t"

type Source interface {
	Start(src chan<- t.Ticks)
}
