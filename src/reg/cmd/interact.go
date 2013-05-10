package cmd

import (
	"bufio"
	"log"
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

	cmdin, err := cmdc.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	var cmdout *bufio.Reader
	if out != nil {
		cmdo, err := cmdc.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		cmdout = bufio.NewReader(cmdo)
	}

	err = cmdc.Start()
	if err != nil {
		log.Fatal(err)
	}

	for input := range in {
		_, err := cmdin.Write([]byte(strings.Join(input, " ") + "\n"))
		if err != nil {
			log.Fatal(err)
		}
		if out != nil {
			s, err := cmdout.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			out <- s[:len(s)-1]
		}
	}

}
