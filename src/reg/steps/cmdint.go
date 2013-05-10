package steps

import (
	"reg/cmd"
	"reg/t"
)

type stepsource_cmdint struct {
	stepsource_cmd
}

func MakeInteractiveCommandSource(cmd string, sourcetype int) Source {
	return &stepsource_cmdint{stepsource_cmd{cmd: cmd, sourcetype: sourcetype}}
}

func (ss *stepsource_cmdint) Start(src <-chan t.Ticks, prod chan<- t.TicksSteps) {
	stepsource_cmd_process(src, prod, ss.sourcetype, cmd.MakeInteractiveCommand(ss.cmd))
}
