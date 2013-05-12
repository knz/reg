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
		cmdc := exec.Command(shell, "-c", c.cmd)
		if in != nil {
			input := <-in
			cmdc.Args = append(cmdc.Args, input...)
		}
		if out == nil {
			err := cmdc.Run()
			CheckErrIsNil(err, cmdc.Args, ":Run()")
		} else {
			output, err := cmdc.Output()
			CheckErrIsNil(err, cmdc.Args, ":Output()")
			out <- string(output[:len(output)-1])
		}
	}

}
