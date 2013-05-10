package main

import (
	"log"
	"os"
	"reg"
	"reg/act"
	"reg/steps"
	"reg/t"
	"reg/ticks"
	"time"
)

func main() {
	var ts ticks.Source
	var ss steps.Source

	i := 1
	for len(os.Args) > i {
		switch os.Args[i] {
		case "ticks":
			switch os.Args[i+1] {
			case "dummy":
				ts = ticks.MakeDummySource()
				i += 1
			case "time":
				d, err := time.ParseDuration(os.Args[i+2])
				if err != nil {
					log.Fatal(err)
				}
				ts = ticks.MakeTimerSource(d)
				i += 2
			default:
				st := 0
				switch os.Args[i+2] {
				case "monotonic":
					st = t.SRC_MONOTONIC
				case "deltas_only":
					st = t.SRC_DELTAS_ONLY
				default:
					st = t.SRC_INIT_THEN_DELTAS
				}
				ts = ticks.MakeCommandSource(os.Args[i+3], st)
				i += 3
			}
		case "steps":
			switch os.Args[i+1] {
			case "dummy":
				ss = steps.MakeDummySource()
				i += 1
			default:
				st := 0
				switch os.Args[i+2] {
				case "monotonic":
					st = t.SRC_MONOTONIC
				case "deltas_only":
					st = t.SRC_DELTAS_ONLY
				default:
					st = t.SRC_INIT_THEN_DELTAS
				}
				switch os.Args[i+1] {
				case "one":
					ss = steps.MakeCommandSource(os.Args[i+3], st)
				default:
					ss = steps.MakeInteractiveCommandSource(os.Args[i+3], st)
				}
				i += 3
			}
		default:
			i += 1
		}

	}

	//	a := act.MakePrinterActuator(os.Stderr)
	// a := act.MakeDummyActuator()
	// a := act.MakeCommandActuator("echo ACT $0 $@ >/dev/tty")
	a := act.MakeInteractiveCommandActuator("while true; do read a; echo ACT $a >/dev/tty; done")
	d := reg.MakeDomain("default", ts, ss, a)
	d.ThrottleType = reg.ThrottleTicks
	d.ThrottleMinPeriod = 0.01
	d.OutputFile = "/dev/stdout"
	d.AddResource("time", "while true; do read a || break; LANG=C ps -o rss= -p 28403; done")
	d.Start(os.Stdin)
	d.Wait()
}
