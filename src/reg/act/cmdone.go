package act

import (
	"fmt"
	"reg/cmd"
	"reg/t"
)

type actuator_cmd struct {
	cmd string
}

type actuator_cmdone struct {
	actuator_cmd
}

func MakeCommandActuator(cmd string) Actuator {
	return &actuator_cmdone{actuator_cmd{cmd: cmd}}
}

func actuator_cmd_process(src <-chan t.Status, cmdc cmd.Cmd) {
	cmdin := make(chan []string)
	go cmdc.Start(cmdin, nil)

	for a := range src {
		args := make([]string, 7)
		args[0] = a.DomainLabel
		args[1] = fmt.Sprint(a.Ticks)
		args[2] = fmt.Sprint(a.TicksDelta)
		args[3] = fmt.Sprint(a.Steps)
		args[4] = fmt.Sprint(a.StepsDelta)
		args[5] = fmt.Sprint(a.Supply)
		args[6] = fmt.Sprint(a.Delta)
		cmdin <- args
	}
}

func (act *actuator_cmdone) Start(src <-chan t.Status) {
	actuator_cmd_process(src, cmd.MakeOneShotCommand(act.cmd))
}
