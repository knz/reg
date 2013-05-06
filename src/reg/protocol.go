package reg

import ("os/exec"; "log"; "fmt")

func (d *Domain) protocol() {
	cmdc := exec.Command("sh", "-c", d.ProtocolCmd)
	cmdin, err := cmdc.StdinPipe()
	if err != nil {	log.Fatal(err)	}
	err = cmdc.Start()
	if err != nil {	log.Fatal(err)	}

	for a := range d.action {
		_, err := fmt.Fprintln(cmdin, d.resources[a.bin].label, " ", a.currentsupply, " ", a.delta)
		if err != nil {	log.Fatal(err)	}
	}
}
