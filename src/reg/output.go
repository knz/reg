package reg

import ("log"; "fmt"; "syscall")


func (d *Domain) outmgt() {

	nres := len(d.resources)
	stprev := Status{ ticks: 0, steps: 0, usage:make([]StuffSteps, nres) }

	doit := false

	for {
		select {
		case <- d.outready:
		case <- d.statusctl: doit = true; continue
		}
		if !doit {
			<- d.statusctl
		}
		doit = false

		d.query <- true
		st := <- d.status

		if st.ticks == 0 {
			d.out <- fmt.Sprint(d.Label, " -\n")
		} else {
			msg := fmt.Sprint(d.Label, " ",
				st.ticks, " ", st.ticks - stprev.ticks, " ",
				st.steps, " ", st.steps - stprev.steps, " ",
				nres)
			for i := range st.usage {
				msg += fmt.Sprint(" ", d.resources[i].label,
					" ", st.usage[i],
					" ", st.usage[i] - stprev.usage[i])
			}
			d.out <- msg + "\n"
			stprev = st
		}
	}
}

func (d *Domain) output() {
	fd, err := syscall.Open(d.OutputFile, syscall.O_WRONLY | syscall.O_CREAT | syscall.O_TRUNC, 0666)
	if err != nil { log.Fatal(err) }

	set := syscall.FdSet{}
	for {
		set.Bits[fd / 64] = int32(fd) % 64
		syscall.Select(fd + 1, nil, &set, nil, nil);
		d.outready <- true;
		cmd := <- d.out;
		syscall.Write(fd, []byte(cmd))
	}
}
