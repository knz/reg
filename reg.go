package main

import (
	"log"
	"os"
	"reg"
	"reg/ticks"
	"time"
)

func main() {
	var ts ticks.Source
	ts = nil

	if len(os.Args) > 1 {
		if os.Args[1] == "time" {
			d, err := time.ParseDuration(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}
			ts = ticks.MakeTimerSource(d)
		} else {
			st := 0
			switch os.Args[2] {
			case "monotonic":
				st = ticks.TS_MONOTONIC
			case "deltas_only":
				st = ticks.TS_DELTAS_ONLY
			default:
				st = ticks.TS_INIT_THEN_DELTAS
			}
			ts = ticks.MakeCommandSource(os.Args[3], st)
		}
	}
	d := reg.MakeDomain("default", ts)
	d.ThrottleType = reg.ThrottleTicks
	d.ThrottleMinPeriod = 0.01
	d.OutputFile = "/dev/stdout"
	d.StepsCmd = "while true; do read a || break; LANG=C ps -o cputime= -p 28403|tr ':.' '  '| LANG=C awk '{print $1*60+$2+$3/100. }'; done"
	d.AddResource("time", "while true; do read a || break; LANG=C ps -o cputime= -p 28403|tr ':.' '  '|LANG=C awk '{print $1*60+$2+$3/100. }'; done")
	d.Start(os.Stdin)
	d.Wait()
}
