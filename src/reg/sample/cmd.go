// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sample

import (
	"fmt"
	"log"
	"reg/cmd"
	"reg/t"
)

type sampler_cmd struct{ cmd cmd.Cmd }

func MakeCommandSampler(cmd cmd.Cmd) Sampler {
	return &sampler_cmd{cmd}
}

func (s *sampler_cmd) Start(src <-chan t.TicksSteps, prod chan<- t.Sample) {
	cmdin := make(chan []string)
	cmdout := make(chan string)

	go s.cmd.Start(cmdin, cmdout)

	for ts := range src {
		args := make([]string, 2)
		args[0] = fmt.Sprint(ts.Ticks)
		args[1] = fmt.Sprint(ts.Steps)
		cmdin <- args
		output := <-cmdout

		s := t.Sample{Ticks: ts.Ticks, Steps: ts.Steps}

		n, err := fmt.Sscan(output, &s.Usage)
		if n != 1 || err != nil {
			log.Fatal(err)
		}
		prod <- s
	}

}
