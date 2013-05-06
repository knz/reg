package ticks

import "reg/t"

type Source interface {
	Start()
	SetSource(src chan t.Ticks)
}
