package ticks

import ("log"; "os"; "os/exec"; "bufio"; "strconv"; "reg/t")

const (
	TS_MONOTONIC = iota
	TS_INIT_THEN_DELTAS
	TS_DELTAS_ONLY
)

type ticksource_cmd struct {
	ticksource_common
	cmd string
	sourcetype int // TS_* above
}

func (ts *ticksource_cmd) Start() {
	if ts.source == nil { log.Fatal("no source channel connected") }

	shell := os.Getenv("SHELL")
	if shell == "" { shell = "sh" }
	cmdc := exec.Command(shell, "-c", ts.cmd)
	cmdout, err := cmdc.StdoutPipe()
	if err != nil {	log.Fatal(err)	}
	err = cmdc.Start()
	if err != nil {	log.Fatal(err)	}
	reader := bufio.NewReader(cmdout)

	go func() {
		lastval := t.Ticks(-1)
		for {
			tickstr, _ := reader.ReadString('\n')
			v, err := strconv.ParseFloat(tickstr[:len(tickstr)-1], 64)
			if err != nil {	log.Fatal(err)	}
			val := t.Ticks(v)
			if (ts.sourcetype == TS_MONOTONIC) {
				if (lastval < 0) {
					// first value: forward init
					ts.source <- val
				} else {
					// next value: compute delta
					ts.source <- val - lastval
				}
				lastval = val
			} else {
				if (ts.sourcetype == TS_DELTAS_ONLY) {
					// no initial value provided by command, fake one
					ts.source <- t.Ticks(0)
				}
				// forward either init or delta
				ts.source <- val
			}
		}
	}()
}

func MakeCommandSource(cmd string, sourcetype int) Source {
	return &ticksource_cmd{cmd:cmd, sourcetype:sourcetype}
}
