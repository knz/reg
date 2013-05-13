package reg

import . "assert"
import (
	"fmt"
	"reg/t"
	"syscall"
)

func outmgt(flood bool, statusctl <-chan bool,
	status <-chan t.Status, query chan<- bool,
	out chan<- string, outready <-chan bool) {

	doit := false

	for {
		select {
		case <-outready:
		case <-statusctl:
			doit = true
			continue
		}
		if !flood && !doit {
			<-statusctl
		}
		doit = false

		query <- true
		st := <-status

		msg := fmt.Sprint(st.Ticks, st.TicksDelta,
			st.Steps, st.StepsDelta,
			st.Supply, st.Delta, "\n")
		out <- msg
	}
}

func output(fd uintptr, out <-chan string, outready chan<- bool) {

	set := syscall.FdSet{}
	for {
		set.Bits[fd/64] = int32(fd) % 64
		err := syscall.Select(int(fd+1), nil, &set, nil, nil)
		Assert(err == nil, "Select() for write on fd ", fd, ":", err)

		outready <- true
		cmd := <-out
		_, err = syscall.Write(int(fd), []byte(cmd))
		Assert(err == nil, "Write() on fd ", fd, ":", err)
	}
}
