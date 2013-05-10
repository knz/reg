package steps

import (
	"fmt"
	"log"
	"reg/cmd"
	"reg/t"
	"strconv"
)

type stepsource_cmd struct {
	cmd        string
	sourcetype int
}

type stepsource_cmdone struct {
	stepsource_cmd
}

func MakeCommandSource(cmd string, sourcetype int) Source {
	return &stepsource_cmdone{stepsource_cmd{cmd: cmd, sourcetype: sourcetype}}
}

func stepsource_cmd_process(ticks <-chan t.Ticks, src chan<- t.TicksSteps,
	sourcetype int, cmdc cmd.Cmd) {

	cmdin := make(chan []string)
	cmdout := make(chan string)

	go cmdc.Start(cmdin, cmdout)

	if sourcetype == t.SRC_DELTAS_ONLY {
		src <- t.TicksSteps{Ticks: <-ticks, Steps: 0}
	}

	var lastval t.Steps
	for {
		tval := <-ticks

		args := make([]string, 1)
		args[0] = fmt.Sprint(tval)
		cmdin <- args
		output := <-cmdout
		v, err := strconv.ParseFloat(output, 64)
		if err != nil {
			log.Fatal(err)
		}
		sval := t.Steps(v)

		if sourcetype == t.SRC_MONOTONIC {
			tmp := sval - lastval
			lastval = sval
			sval = tmp
		}

		src <- t.TicksSteps{Ticks: tval, Steps: sval}
	}
}

func (ss *stepsource_cmdone) Start(src <-chan t.Ticks, prod chan<- t.TicksSteps) {
	stepsource_cmd_process(src, prod, ss.sourcetype, cmd.MakeOneShotCommand(ss.cmd))
}
