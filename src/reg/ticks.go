package reg

import ("os/exec"; "log"; "strconv"; "bufio"; "time"; "os")

const (
	TS_MONOTONIC = iota
	TS_INIT_THEN_DELTAS
	TS_DELTAS_ONLY
)

type ticksource_common struct {
	source chan Ticks
}


type ticksource_cmd struct {
	ticksource_common
	cmd string
	sourcetype int // TS_* above
}

type ticksource_timer struct {
	ticksource_common
	period time.Duration
}

func (ts *ticksource_common) GetSource() chan Ticks {
	return ts.source
}
func (ts *ticksource_cmd) Start() {
	shell := os.Getenv("SHELL")
	if shell == "" { shell = "sh" }
	cmdc := exec.Command(shell, "-c", ts.cmd)
	cmdout, err := cmdc.StdoutPipe()
	if err != nil {	log.Fatal(err)	}
	err = cmdc.Start()
	if err != nil {	log.Fatal(err)	}
	reader := bufio.NewReader(cmdout)

	go func() {
		lastval := Ticks(-1)
		for {
			tickstr, _ := reader.ReadString('\n')
			v, err := strconv.ParseFloat(tickstr[:len(tickstr)-1], 64)
			if err != nil {	log.Fatal(err)	}
			t := Ticks(v)
			if (ts.sourcetype == TS_MONOTONIC) {
				if (lastval < 0) {
					// first value: forward init
					ts.source <- t
				} else {
					// next value: compute delta
					ts.source <- t - lastval
				}
				lastval = t
			} else {
				if (ts.sourcetype == TS_DELTAS_ONLY) {
					// no initial value provided by command, fake one
					ts.source <- Ticks(0)
				}
				// forward either init or delta
				ts.source <- t
			}
		}
	}()
}

func MakeCommandSource(cmd string, sourcetype int) TickSource {
	v := &ticksource_cmd{cmd:cmd, sourcetype:sourcetype}
	v.source = make(chan Ticks)
	return v
}


func MakeTimerSource(period time.Duration) TickSource {
	v := &ticksource_timer{period:period}
	v.source = make(chan Ticks)
	return v
}

func (ts *ticksource_timer) Start() {
	go func() {
		ts.source <- Ticks(0) // initial
		for {
			time.Sleep(ts.period)
			ts.source <- Ticks(1) // delta
		}
	}()
}

func (d *Domain) ticksource() {
	ticksext := d.TickSource.GetSource()
	d.tickssrc <- <- ticksext // forward init
	for {
		// forward deltas from either external tick source
		// or input stream
		t := Ticks(0)
		select {
		case t = <- d.ticksctl:
		case t = <- ticksext:
		}
		if t > 0 {
			d.tickssrc <- t
		}
	}
}
