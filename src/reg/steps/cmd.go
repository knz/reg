package steps

import . "assert"
import (
	"fmt"
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

	if (ss.sourcetype&t.SRC_Z != 0) || (ss.sourcetype&t.SRC_O == 0) {
		// we produce zero as origin in this case.
		// If the step function also produces an origin, just wait for it and then ignore.
		tsrc := <-src
		if ss.sourcetype&t.SRC_O != 0 {
			args := make([]string, 1)
			args[0] = fmt.Sprint(tsrc)
			cmdin <- args
			<-cmdout // drop it
		}
		prod <- t.TicksSteps{Ticks: tsrc, Steps: 0}
	}

	var lastval t.Steps
	for tval := range src {

		args := make([]string, 1)
		args[0] = fmt.Sprint(tval)
		cmdin <- args
		output := <-cmdout
		v, err := strconv.ParseFloat(output, 64)
		Assert(err == nil, "parsing steps", ":", err)

		sval := t.Steps(v)

		if ss.sourcetype&t.SRC_M != 0 {
			tmp := sval - lastval
			lastval = sval
			sval = tmp
		}

		prod <- t.TicksSteps{Ticks: tval, Steps: sval}
	}
}
