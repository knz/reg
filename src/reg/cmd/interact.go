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
		CheckErrIsNil(err, cmdc.Args, ":StdinPipe()")
		cmdin = cmdi
	}

	var cmdout *bufio.Reader
	if out != nil {
		cmdo, err := cmdc.StdoutPipe()
		CheckErrIsNil(err, cmdc.Args, ":StdoutPipe()")
		cmdout = bufio.NewReader(cmdo)
	}

	err := cmdc.Start()
	CheckErrIsNil(err, cmdc.Args, ":Start()")

	for {
		if in != nil {
			input := <-in
			_, err := cmdin.Write([]byte(strings.Join(input, " ") + "\n"))
			CheckErrIsNil(err, cmdc.Args, ":Write()")
		}
		if out != nil {
			s, err := cmdout.ReadString('\n')
			CheckErrIsNil(err, cmdc.Args, ":ReadString()")
			out <- s[:len(s)-1]
		}
	}

}
