package ticks

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"reg/t"
	"strconv"
)

type ticksource_cmd struct {
	cmd        string
	sourcetype int
}

func (ts *ticksource_cmd) Start(prod chan<- t.Ticks) {

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	cmdc := exec.Command(shell, "-c", ts.cmd)
	cmdout, err := cmdc.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmdc.Start()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(cmdout)

	lastval := t.Ticks(0)
	if ts.sourcetype == t.SRC_DELTAS_ONLY {
		// no initial value provided by command, fake one
		prod <- t.Ticks(0)
	}

	for {
		tickstr, _ := reader.ReadString('\n')
		v, err := strconv.ParseFloat(tickstr[:len(tickstr)-1], 64)
		if err != nil {
			log.Fatal(err)
		}
		val := t.Ticks(v)

		if ts.sourcetype == t.SRC_MONOTONIC {
			tmp := val - lastval
			lastval = val
			val = tmp
		}

		prod <- val
	}
}

func MakeCommandSource(cmd string, sourcetype int) Source {
	return &ticksource_cmd{cmd: cmd, sourcetype: sourcetype}
}
