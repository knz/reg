package ticks

import ("log"; "reg/t")

type ticksource_dummy struct { ticksource_common }

func (ts *ticksource_dummy) Start() {
	if ts.source == nil { log.Fatal("no source channel connected") }
	go func() {
		ts.source <- t.Ticks(0) // init,
		// then nothing
	}()
}

func MakeDummySource() Source {
	return &ticksource_dummy{}
}
