package main

import ("os"; "reg"; "time"; "log")

func main() {
	var ts reg.TickSource
	ts = nil

	if len(os.Args) > 1 {
		if os.Args[1] == "time" {
			d, err := time.ParseDuration(os.Args[2])
			if err != nil {	log.Fatal(err)}
			ts = reg.MakeTimerSource(d)
		} else {
			st := 0
			switch os.Args[2] {
			case "monotonic": st = reg.TS_MONOTONIC
			case "deltas_only": st = reg.TS_DELTAS_ONLY
			default: st = reg.TS_INIT_THEN_DELTAS
			}
			ts = reg.MakeCommandSource(os.Args[3], st)
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
