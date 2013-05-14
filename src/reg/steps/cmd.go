// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package steps

import . "assert"
import (
	"fmt"
	"reg/cmd"
	"reg/t"
	"strconv"
)

type stepsource_cmd struct {
	cmd        cmd.Cmd
	sourcetype int
}

func MakeCommandSource(cmd cmd.Cmd, sourcetype int) Source {
	return &stepsource_cmd{cmd, sourcetype}
}

func (ss *stepsource_cmd) Start(src <-chan t.Ticks, prod chan<- t.TicksSteps) {

	cmdin := make(chan []string)
	cmdout := make(chan string)

	go ss.cmd.Start(cmdin, cmdout)

	if ss.sourcetype&t.SRC_O != 0 {
		// the main loop below skips zero step increases.
		// however we must guarantee at least one step event is issued
		// as origin.
		tval := <-src
		args := make([]string, 1)
		args[0] = fmt.Sprint(tval)
		cmdin <- args

		sval := t.Steps(0)
		if ss.sourcetype&t.SRC_Z != 0 {
			// we produce zero as origin in this case
			<-cmdout
		} else {
			sval_str := <-cmdout
			v, err := strconv.ParseFloat(sval_str, 64)
			Assert(err == nil, "parsing steps", ":", err)
			sval = t.Steps(v)
		}
		prod <- t.TicksSteps{Ticks: tval, Steps: sval}

	} else {
		if ss.sourcetype&t.SRC_Z != 0 {
			// we produce zero as origin in this case
			prod <- t.TicksSteps{Ticks: <-src, Steps: 0}
		}
	}

	var lastval t.Steps
	for tval := range src {

		args := make([]string, 1)
		args[0] = fmt.Sprint(tval)
		cmdin <- args
		output := <-cmdout
		v, err := strconv.ParseFloat(output, 64)
		Assert(err == nil, "parsing steps", ":", err)

		sval := t.Steps(v)

		if ss.sourcetype&t.SRC_M != 0 {
			tmp := sval - lastval
			lastval = sval
			sval = tmp
		}

		if sval > 0 {
			prod <- t.TicksSteps{Ticks: tval, Steps: sval}
		}
	}
}
