package act

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type actuator_cmd struct {
	cmd string
}

type actuator_cmdone struct {
	actuator_common
	actuator_cmd
}

func MakeCommandActuator(cmd string) Actuator {
	return &actuator_cmdone{actuator_common{}, actuator_cmd{cmd: cmd}}
}

func (act *actuator_cmdone) Start() {
	act.Check()

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}

	go func() {
		for a := range act.source {
			args := make([]string, 7+2*len(a.Supply))
			args[0] = "-c"
			args[1] = act.cmd
			args[2] = a.DomainLabel
			args[3] = fmt.Sprint(a.Ticks)
			args[4] = fmt.Sprint(a.TicksDelta)
			args[5] = fmt.Sprint(a.Steps)
			args[6] = fmt.Sprint(a.StepsDelta)
			for i, s := range a.Supply {
				args[7+2*i] = fmt.Sprint(s)
				args[7+2*i+1] = fmt.Sprint(a.Delta[i])
			}
			cmdc := exec.Command(shell, args...)
			err := cmdc.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}
