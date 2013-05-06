package reg

import ("io"; "bufio"; "strconv"; "strings")

func (d *Domain) readlines(input io.Reader) {
	reader := bufio.NewReader(input)
	for {
		cmdstr, err := reader.ReadString('\n')
		if err == io.EOF { d.inputdone <- true; break }
		d.input <- cmdstr[:len(cmdstr)-1]
	}
}

func (d *Domain) parse() {
	for cmd := range d.input {
		cmdargs := strings.Split(cmd, " ");

		switch cmdargs[0] {
		case ".":
			v, _ := strconv.ParseFloat(cmdargs[1], 64)
			d.ticksctl <- Ticks(v)
		case "+":
			b, _ := strconv.ParseInt(cmdargs[1], 0, 0)
			v, _ := strconv.ParseFloat(cmdargs[2], 64)
			d.supplycmd <- SupplyCmd{bin : int(b), supply : StuffSteps(v)}
		case "?":
			d.statusctl <- true
		}
	}
}
