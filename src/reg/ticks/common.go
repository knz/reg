package ticks

import (
	"log"
	"reg/t"
)

type ticksource_common struct {
	source chan t.Ticks
}

func (ts *ticksource_common) SetSource(src chan t.Ticks) {
	ts.source = src
}

func (ts *ticksource_common) Check() {
	if ts.source == nil {
		log.Fatal("no source channel connected")
	}
}
