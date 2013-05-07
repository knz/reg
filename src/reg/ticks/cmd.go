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
	ticksource_common
	cmd        string
	sourcetype int
}

func (ts *ticksource_cmd) Start() {
	ts.Check()

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

	go func() {
		lastval := t.Ticks(0)
		if ts.sourcetype == t.SRC_DELTAS_ONLY {
			// no initial value provided by command, fake one
			ts.source <- t.Ticks(0)
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

			ts.source <- val
		}
	}()
}

func MakeCommandSource(cmd string, sourcetype int) Source {
	return &ticksource_cmd{cmd: cmd, sourcetype: sourcetype}
}
