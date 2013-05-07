package steps

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reg/t"
	"strconv"
)

type stepsource_cmd struct {
	cmd        string
	sourcetype int
}

type stepsource_cmdone struct {
	stepsource_common
	stepsource_cmd
}

func MakeCommandSource(cmd string, sourcetype int) Source {
	return &stepsource_cmdone{stepsource_common{}, stepsource_cmd{cmd: cmd, sourcetype: sourcetype}}
}

func stepsource_cmd_process(ticks chan t.Ticks, src chan t.TicksSteps,
	sourcetype int, interact func(ticks t.Ticks) t.Steps) {

	if sourcetype == t.SRC_DELTAS_ONLY {
		src <- t.TicksSteps{Ticks: <-ticks, Steps: 0}
	}

	var lastval t.Steps
	for {
		tval := <-ticks
		sval := interact(tval)

		if sourcetype == t.SRC_MONOTONIC {
			tmp := sval - lastval
			lastval = sval
			sval = tmp
		}

		src <- t.TicksSteps{Ticks: tval, Steps: sval}
	}
}

func (ss *stepsource_cmdone) Start() {
	ss.Check()

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}

	go stepsource_cmd_process(ss.ticks, ss.source, ss.sourcetype, func(ticks t.Ticks) t.Steps {
		cmdc := exec.Command(shell, "-c", ss.cmd, fmt.Sprint(ticks))
		output, err := cmdc.Output()
		if err != nil {
			log.Fatal(err)
		}
		v, err := strconv.ParseFloat(string(output[:len(output)-1]), 64)
		if err != nil {
			log.Fatal(err)
		}
		return t.Steps(v)
	})
}
