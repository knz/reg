package act

import (
	"fmt"
	"io"
)

type actuator_printer struct {
	actuator_common
	out io.Writer
}

func MakePrinterActuator(out io.Writer) Actuator {
	return &actuator_printer{actuator_common{}, out}
}

func (act *actuator_printer) Start() {
	go func() {
		for action := range act.source {
			s := fmt.Sprint(action.DomainLabel, " ",
				action.TicksDelta, " ",
				action.StepsDelta, " ",
				len(action.Supply))
			for i, v := range action.Supply {
				s += fmt.Sprint(" ",
					v, " ", action.Delta[i])
			}
			s += "\n"
			act.out.Write([]byte(s))
		}
	}()
}
