package reg

import (
	"fmt"
	"log"
	"reg/t"
	"syscall"
)

func (d *Domain) outmgt(statusctl <-chan bool, status <-chan t.Status, query chan<- bool, out chan<- string, outready <-chan bool) {

	doit := false

	for {
		select {
		case <-outready:
		case <-statusctl:
			doit = true
			continue
		}
		if !doit {
			<-statusctl
		}
		doit = false

		query <- true
		st := <-status

		msg := fmt.Sprint(d.Label, " ",
			st.Ticks, st.TicksDelta,
			st.Steps, st.StepsDelta,
			st.Supply, st.Delta, "\n")
		out <- msg
	}
}

func (d *Domain) output(out <-chan string, outready chan<- bool) {
	fd, err := syscall.Open(d.OutputFile, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	set := syscall.FdSet{}
	for {
		set.Bits[fd/64] = int32(fd) % 64
		syscall.Select(fd+1, nil, &set, nil, nil)
		outready <- true
		cmd := <-out
		syscall.Write(fd, []byte(cmd))
	}
}
