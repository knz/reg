package reg

import (
	"fmt"
	"log"
	"syscall"
)

func (d *Domain) outmgt() {

	doit := false

	for {
		select {
		case <-d.outready:
		case <-d.statusctl:
			doit = true
			continue
		}
		if !doit {
			<-d.statusctl
		}
		doit = false

		d.query <- true
		st := <-d.status

		msg := fmt.Sprint(d.Label, " ",
			st.Ticks, " ", st.TicksDelta, " ",
			st.Steps, " ", st.StepsDelta)
		for i, v := range st.Supply {
			msg += fmt.Sprint(" ", v, " ", st.Delta[i])
		}
		d.out <- msg + "\n"
	}
}

func (d *Domain) output() {
	fd, err := syscall.Open(d.OutputFile, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	set := syscall.FdSet{}
	for {
		set.Bits[fd/64] = int32(fd) % 64
		syscall.Select(fd+1, nil, &set, nil, nil)
		d.outready <- true
		cmd := <-d.out
		syscall.Write(fd, []byte(cmd))
	}
}
