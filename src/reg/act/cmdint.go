package act

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type actuator_cmdint struct {
	actuator_common
	actuator_cmd
}

func MakeInteractiveCommandActuator(cmd string) Actuator {
	return &actuator_cmdint{actuator_common{}, actuator_cmd{cmd: cmd}}
}

func (act *actuator_cmdint) Start() {
	act.Check()

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	cmdc := exec.Command(shell, "-c", act.cmd)
	cmdin, err := cmdc.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmdc.Start()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for a := range act.source {
			msg := fmt.Sprint(a.DomainLabel, " ",
				a.Ticks, " ", a.TicksDelta, " ",
				a.Steps, " ", a.StepsDelta)
			for i, v := range a.Supply {
				msg += fmt.Sprint(" ", v, " ", a.Delta[i])
			}
			msg += "\n"
			_, err := cmdin.Write([]byte(msg))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}
