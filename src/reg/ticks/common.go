package ticks

type ticksource_common struct {
	source chan Ticks
}

func (ts *ticksource_common) SetSource(src chan Ticks)  {
	ts.source = src
}
