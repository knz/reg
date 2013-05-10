package act

import (
	"fmt"
	"io"
	"reg/t"
)

type actuator_printer struct {
	out io.Writer
}

func MakePrinterActuator(out io.Writer) Actuator {
	return &actuator_printer{out}
}

func (act *actuator_printer) Start(src <-chan t.Status) {
	for action := range src {
		s := fmt.Sprint(action.DomainLabel, " ",
			action.Ticks, action.TicksDelta,
			action.Steps, action.StepsDelta,
			action.Supply, action.Delta, "\n")
		act.out.Write([]byte(s))
	}
}
