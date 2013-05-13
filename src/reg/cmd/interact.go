// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import . "assert"
import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

type cmd_interactive struct {
	cmd string
}

func MakeInteractiveCommand(cmd string) Cmd {
	return &cmd_interactive{cmd}
}

func (c *cmd_interactive) Start(in <-chan []string, out chan<- string) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	cmdc := exec.Command(shell, "-c", c.cmd)

	var cmdin io.WriteCloser
	if in != nil {
		cmdi, err := cmdc.StdinPipe()
		Assert(err == nil, cmdc.Args, ":StdinPipe()", ":", err)
		cmdin = cmdi
	}

	var cmdout *bufio.Reader
	if out != nil {
		cmdo, err := cmdc.StdoutPipe()
		Assert(err == nil, cmdc.Args, ":StdoutPipe()", ":", err)
		cmdout = bufio.NewReader(cmdo)
	}

	err := cmdc.Start()
	Assert(err == nil, cmdc.Args, ":Start()", ":", err)

	for {
		if in != nil {
			input := <-in
			_, err := cmdin.Write([]byte(strings.Join(input, " ") + "\n"))
			Assert(err == nil, cmdc.Args, ":Write()", ":", err)
		}
		if out != nil {
			s, err := cmdout.ReadString('\n')
			Assert(err == nil, cmdc.Args, ":ReadString()", ":", err)
			out <- s[:len(s)-1]
		}
	}

}
