package act

import (
	"fmt"
	"reg/cmd"
	"reg/t"
)

type actuator_cmd struct {
	cmd cmd.Cmd
}

func MakeCommandActuator(cmd cmd.Cmd) Actuator {
	return &actuator_cmd{cmd}
}

func (act *actuator_cmd) Start(src <-chan t.Status) {
	cmdin := make(chan []string)
	go act.cmd.Start(cmdin, nil)

	for a := range src {
		args := make([]string, 6)
		args[0] = fmt.Sprint(a.Ticks)
		args[1] = fmt.Sprint(a.TicksDelta)
		args[2] = fmt.Sprint(a.Steps)
		args[3] = fmt.Sprint(a.StepsDelta)
		args[4] = fmt.Sprint(a.Supply)
		args[5] = fmt.Sprint(a.Delta)
		cmdin <- args
	}
}
