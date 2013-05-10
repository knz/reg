package cmd

import (
	"log"
	"os"
	"os/exec"
)

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

	for input := range in {

		cmdc := exec.Command(shell, "-c", c.cmd)
		cmdc.Args = append(cmdc.Args, input...)
		if out == nil {
			err := cmdc.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			output, err := cmdc.Output()
			if err != nil {
				log.Fatal(err)
			}
			out <- string(output[:len(output)-1])
		}
	}

}
