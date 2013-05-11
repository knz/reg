package steps

import (
	"fmt"
	"log"
	"reg/cmd"
	"reg/t"
	"strconv"
)

type stepsource_cmd struct {
	cmd        cmd.Cmd
	sourcetype int
}

func MakeCommandSource(cmd cmd.Cmd, sourcetype int) Source {
	return &stepsource_cmd{cmd, sourcetype}
}

func (ss *stepsource_cmd) Start(src <-chan t.Ticks, prod chan<- t.TicksSteps) {

	cmdin := make(chan []string)
	cmdout := make(chan string)

	go ss.cmd.Start(cmdin, cmdout)

	if ss.sourcetype == t.SRC_DELTAS_ONLY {
		prod <- t.TicksSteps{Ticks: <-src, Steps: 0}
	}

	var lastval t.Steps
	for tval := range src {

		args := make([]string, 1)
		args[0] = fmt.Sprint(tval)
		cmdin <- args
		output := <-cmdout
		v, err := strconv.ParseFloat(output, 64)
		if err != nil {
			log.Fatal(err)
		}
		sval := t.Steps(v)

		if ss.sourcetype == t.SRC_MONOTONIC {
			tmp := sval - lastval
			lastval = sval
			sval = tmp
		}

		prod <- t.TicksSteps{Ticks: tval, Steps: sval}
	}
}
