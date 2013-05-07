package reg

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"reg/t"
)

func (d *Domain) sample() {

	nres := len(d.resources)

	type ResourceCmd struct {
		cmd *exec.Cmd
		out io.ReadCloser
		in  io.WriteCloser
	}
	cmds := make([]ResourceCmd, nres)
	for i := range cmds {
		cmds[i].cmd = exec.Command("sh", "-c", d.resources[i].cmd)
		cmdout, err := cmds[i].cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		cmds[i].out = cmdout
		cmds[i].in, err = cmds[i].cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		err = cmds[i].cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}

	for st := range d.tickssteps {
		values := make([]t.Stuff, nres)

		for i := range values {
			_, err := fmt.Fprintln(cmds[i].in, st.ticks, ' ', st.steps)
			if err != nil {
				log.Fatal(err)
			}
		}
		for i := range values {
			n, err := fmt.Fscanln(cmds[i].out, &values[i])
			if err != nil || n != 1 {
				log.Fatal(err)
			}
		}

		d.measure <- Sample{ticks: st.ticks, steps: st.steps, usage: values}
	}
}
