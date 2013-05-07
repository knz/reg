package steps

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reg/t"
)

type stepsource_cmdint struct {
	stepsource_common
	stepsource_cmd
}

func MakeInteractiveCommandSource(cmd string, sourcetype int) Source {
	return &stepsource_cmdint{stepsource_common{}, stepsource_cmd{cmd: cmd, sourcetype: sourcetype}}
}

func (ss *stepsource_cmdint) Start() {
	ss.Check()

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	cmdc := exec.Command(shell, "-c", ss.cmd)
	cmdout, err := cmdc.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmdin, err := cmdc.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmdc.Start()
	if err != nil {
		log.Fatal(err)
	}

	go stepsource_cmd_process(ss.ticks, ss.source, ss.sourcetype, func(ticks t.Ticks) t.Steps {
		var s t.Steps
		n, err := fmt.Fprintln(cmdin, ticks)
		if err != nil {
			log.Panic(err)
		}
		n, err = fmt.Fscanln(cmdout, &s)
		if err != nil || n != 1 {
			log.Panic(err)
		}

		return s
	})
}
