package ticks
import "reg/t"
type ticksource_common struct {
	source chan t.Ticks
}

func (ts *ticksource_common) SetSource(src chan t.Ticks)  {
	ts.source = src
}
