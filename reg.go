package main

import . "assert"
import (
	"os"
	"reg"
	"reg/act"
	"reg/cmd"
	"reg/sample"
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
				CheckErrIsNil(err, "parsing command-line argument")
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
				switch os.Args[i+1] {
				case "one":
					ts = ticks.MakeCommandSource(cmd.MakeOneShotCommand(os.Args[i+3]), st)
				default:
					ts = ticks.MakeCommandSource(cmd.MakeInteractiveCommand(os.Args[i+3]), st)
				}
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
					ss = steps.MakeCommandSource(cmd.MakeOneShotCommand(os.Args[i+3]), st)
				default:
					ss = steps.MakeCommandSource(cmd.MakeInteractiveCommand(os.Args[i+3]), st)
				}
				i += 3
			}
		default:
			i += 1
		}

	}

	/*
		var fin, fout os.File
		var err Error

		if inputfile == nil || inputfile == "-" {
			fin = os.Stdin
		} else {
			fin, err = os.Open(inputfile)
			if err != nil {
				log.Fatal(err)
			}
		}

		if outputfile == nil || outputfile == "-" {
			fout = os.Stdout
		} else {
			fout, err := os.OpenFile(outputfile, os.O_WRONLY|os.O_CREAT|os.O_TRUNC, 0666)
			if err != nil {
				log.Fatal(err)
			}
		}

	*/
	//	a := act.MakePrinterActuator(os.Stderr)
	// a := act.MakeDummyActuator()
	// a := act.MakeCommandActuator("echo ACT $0 $@ >/dev/tty")
	a := act.MakeCommandActuator(cmd.MakeInteractiveCommand("while true; do read a || break; echo ACT $a >/dev/tty; done"))

	s := sample.MakeCommandSampler(cmd.MakeOneShotCommand("LANG=C ps -o rss= -p 99298"))
	d := reg.MakeDomain(ts, ss, a, s)
	d.Start(os.Stdin, os.Stdout, reg.OUTPUT_FLOOD, 1, true)
	d.Wait()
}
