// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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
		s := fmt.Sprint(action.Ticks, action.TicksDelta,
			action.Steps, action.StepsDelta,
			action.Supply, action.Delta, "\n")
		act.out.Write([]byte(s))
	}
}
