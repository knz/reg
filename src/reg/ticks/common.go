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

func TeeTicks(dst1 chan t.Ticks, dst2 chan t.Ticks, src chan t.Ticks) {
	go func() {
		for a := range src {
			dst1 <- a
			dst2 <- a
		}
	}()
}
