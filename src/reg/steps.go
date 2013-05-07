package reg

import (
	"fmt"
	"log"
	"os/exec"
	"reg/t"
)

func (d *Domain) stepsource() {
	cmdc := exec.Command("sh", "-c", d.StepsCmd)
	cmdout, err := cmdc.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmdin, err := cmdc.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmdc.Start()
	if err != nil {
		log.Fatal(err)
	}

	s_prev := t.Steps(-1)
	for ticks := range d.ticksin {
		n, err := fmt.Fprintln(cmdin, ticks)
		if err != nil {
			log.Fatal("stepsource.cmdin ", err)
		}

		s := t.Steps(0)
		n, err = fmt.Fscanln(cmdout, &s)
		if err != nil || n != 1 {
			log.Fatal("stepsource.cmdout ", err)
		}

		if s_prev >= 0 {
			s_delta := s - s_prev
			if s_delta > 0 {
				d.tickssteps <- TicksSteps{ticks: ticks, steps: s_delta}
				d.stepsper <- s_delta
			}
		}
		s_prev = s
	}
}
