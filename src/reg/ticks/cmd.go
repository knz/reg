package ticks

import (
	"reg/cmd"
	"reg/t"
	"strconv"
)

import . "assert"

type ticksource_cmd struct {
	cmd        cmd.Cmd
	sourcetype int
}

func MakeCommandSource(cmd cmd.Cmd, sourcetype int) Source {
	return &ticksource_cmd{cmd, sourcetype}
}

func (ts *ticksource_cmd) Start(prod chan<- t.Ticks) {

	cmdout := make(chan string)
	go ts.cmd.Start(nil, cmdout)

	lastval := t.Ticks(0)

	if (ts.sourcetype&t.SRC_Z != 0) || (ts.sourcetype&t.SRC_O == 0) {
		// we produce zero as origin in this case.
		// If the command also produces an origin, just wait for it and then ignore.
		if ts.sourcetype&t.SRC_O != 0 {
			<-cmdout
		}
		// then emit the zero origin.
		prod <- t.Ticks(0)
	}

	for {
		tickstr := <-cmdout

		v, err := strconv.ParseFloat(tickstr, 64)
		CheckErrIsNil(err, "parsing ticks")

		val := t.Ticks(v)

		if ts.sourcetype&t.SRC_M != 0 {
			tmp := val - lastval
			lastval = val
			val = tmp
		}

		prod <- val
	}
}
