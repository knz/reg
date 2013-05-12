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
	if ts.sourcetype == t.SRC_DELTAS_ONLY {
		// no initial value provided by command, fake one
		prod <- t.Ticks(0)
	}

	for {
		tickstr := <-cmdout

		v, err := strconv.ParseFloat(tickstr, 64)
		CheckErrIsNil(err, "parsing ticks")

		val := t.Ticks(v)

		if ts.sourcetype == t.SRC_MONOTONIC {
			tmp := val - lastval
			lastval = val
			val = tmp
		}

		prod <- val
	}
}
