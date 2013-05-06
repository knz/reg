package ticks

type Ticks float64

type Source interface {
	Start()
	SetSource(src chan Ticks)
}
