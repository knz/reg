// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"os"
	"os/exec"
)

import . "assert"

type cmd_oneshot struct {
	cmd string
}

func MakeOneShotCommand(cmd string) Cmd {
	return &cmd_oneshot{cmd}
}

func (c *cmd_oneshot) Start(in <-chan []string, out chan<- string) {

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}

	for {
		cmdc := exec.Command(shell, "-c", c.cmd, "reg")
		if in != nil {
			input := <-in
			cmdc.Args = append(cmdc.Args, input...)
		}
		if out == nil {
			err := cmdc.Run()
			Assert(err == nil, cmdc.Args, ":Run()", ":", err)
		} else {
			output, err := cmdc.Output()
			Assert(err == nil, cmdc.Args, ":Output()", ":", err)
			out <- string(output[:len(output)-1])
		}
	}

}
