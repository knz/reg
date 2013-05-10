package act

import (
	"reg/cmd"
	"reg/t"
)

type actuator_cmdint struct {
	actuator_cmd
}

func MakeInteractiveCommandActuator(cmd string) Actuator {
	return &actuator_cmdint{actuator_cmd{cmd: cmd}}
}

func (act *actuator_cmdint) Start(src <-chan t.Status) {
	actuator_cmd_process(src, cmd.MakeInteractiveCommand(act.cmd))
}
