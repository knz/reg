package main

import . "assert"
import (
	"fmt"
	"getopt"
	"os"
	"reg"
	"reg/act"
	"reg/cmd"
	"reg/sample"
	"reg/steps"
	"reg/t"
	"reg/ticks"
	"strconv"
	"strings"
	"time"
)

func flagconv(flags string) (int, bool) {
	v := 0
	for _, c := range flags {
		switch c {
		case 'd':
			v |= t.SRC_D
		case 'm':
			v |= t.SRC_M
		case 'o':
			v |= t.SRC_O
		case 'z':
			v |= t.SRC_Z
		case 0: // end-of-string, do nothing
		default:
			return 0, false
		}
	}
	return v, true
}

func split(arg string) (string, string, string) {
	wa := strings.SplitN(arg, ":", 2)
	wf := strings.SplitN(wa[0], "/", 2)
	ty := wf[0]
	var a, fl string
	if len(wa) > 1 {
		a = wa[1]
	}
	if len(wf) > 1 {
		fl = wf[1]
	}
	return ty, fl, a
}

func main() {
	ifname := getopt.StringLong("input", 'i', "-", "Set input stream for supplies (default '-', stdin)", "FILE")
	ofname := getopt.StringLong("output", 'o', "-", "Set output stream for reports (default '-', stdout)", "FILE")
	tspec := getopt.StringLong("ticks", 't', "", "Set tick generator (produces tick events)", "SPEC")
	sspec := getopt.StringLong("steps", 's', "", "Set progress indicator (measures steps)", "SPEC")
	mspec := getopt.StringLong("monitor", 'm', "", "Set resource monitor (measures stuff)", "SPEC")
	aspec := getopt.StringLong("actuator", 'a', "", "Set actuator (triggered upon supply exhaustion)", "SPEC")
	gran := getopt.StringLong("granularity", 'g', "0", "Force tick granularity (default 0, disabled)", "N")
	thr := getopt.StringLong("periodic-output", 'p', "none", "Configure periodic output (default none)", "PER")

	getopt.SetParameters("")
	getopt.Parse()

	/*** -i / -o ***/

	var fin, fout *os.File
	var err error

	if *ifname == "-" {
		fin = os.Stdin
	} else {
		fin, err = os.Open(*ifname)
		CheckErrIsNil(err, "-i")
	}

	if *ofname == "-" {
		fout = os.Stdout
	} else {
		fout, err = os.OpenFile(*ofname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		CheckErrIsNil(err, "-o")
	}

	/*** -g ***/

	g, err := strconv.ParseFloat(*gran, 64)
	CheckErrIsNil(err, "-g")

	/*** -p ***/

	throttle_type := reg.OUTPUT_EXPLICIT_ONLY
	throttle_period := float64(0)

	spec_type, spec_flags, spec_arg := split(*thr)
	switch spec_type {
	case "none":
		Assert(spec_flags == "" && spec_arg == "", "-p none does not accept flags/argument")
		throttle_type = reg.OUTPUT_EXPLICIT_ONLY
	case "flood":
		Assert(spec_flags == "" && spec_arg == "", "-p flood does not accept flags/argument")
		throttle_type = reg.OUTPUT_FLOOD
	case "ticks":
		throttle_type = reg.OUTPUT_THROTTLE_TICKS
		throttle_period, err = strconv.ParseFloat(spec_arg, 64)
		CheckErrIsNil(err, "-p ticks")
		Assert(spec_flags == "", "-p ticks does not accept flags")
	case "steps":
		throttle_type = reg.OUTPUT_THROTTLE_STEPS
		throttle_period, err = strconv.ParseFloat(spec_arg, 64)
		CheckErrIsNil(err, "-p steps")
		Assert(spec_flags == "", "-p steps does not accept flags")
	default:
		Assert(false, "invalid syntax for -p")
	}

	/*** -a ***/

	var a act.Actuator
	spec_type, spec_flags, spec_arg = split(*aspec)
	Assert(spec_flags == "", "-a does not accept flags")
	switch spec_type {
	case "dummy":
		Assert(spec_arg == "", "-a dummy does not accept argument")
		a = act.MakeDummyActuator()
	case "print":
		var af *os.File
		if spec_arg == "-" {
			af = os.Stdout
		} else {
			af, err = os.OpenFile(*ofname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			CheckErrIsNil(err, "-a")
		}
		a = act.MakePrinterActuator(af)
	case "proc":
		a = act.MakeCommandActuator(cmd.MakeInteractiveCommand(spec_arg))
	case "cmd":
		a = act.MakeCommandActuator(cmd.MakeOneShotCommand(spec_arg))
	default:
		Assert(false, "invalid -a, or -a not specified")
	}

	/*** -t ***/

	var ts ticks.Source
	spec_type, spec_flags, spec_arg = split(*tspec)
	switch spec_type {
	case "dummy":
		Assert(spec_flags == "" && spec_arg == "", "-t dummy does not accept flags/argument")
		ts = ticks.MakeDummySource()
	case "time", "ptime", "cmd", "proc":
		stype, ok := flagconv(spec_flags)
		fmt.Println("HELLO", stype, ok, spec_flags)
		Assert(ok, "invalid flags for -t")
		switch spec_type {
		case "time", "ptime":
			d, err := time.ParseDuration(spec_arg)
			CheckErrIsNil(err, "-t")
			ts = ticks.MakeTimerSource(d, stype, spec_type == "ptype")
		case "cmd":
			ts = ticks.MakeCommandSource(cmd.MakeOneShotCommand(spec_arg), stype)
		case "proc":
			ts = ticks.MakeCommandSource(cmd.MakeInteractiveCommand(spec_arg), stype)
		}
	default:
		Assert(false, "invalid -t, or -t not specified")
	}

	/*** -s ***/

	var ss steps.Source
	spec_type, spec_flags, spec_arg = split(*sspec)
	switch spec_type {
	case "dummy":
		Assert(spec_flags == "" && spec_arg == "", "-s dummy does not accept flags/argument")
		ss = steps.MakeDummySource()
	case "cmd", "proc":
		stype, ok := flagconv(spec_flags)
		Assert(ok, "invalid flags for -s")
		switch spec_type {
		case "cmd":
			ss = steps.MakeCommandSource(cmd.MakeOneShotCommand(spec_arg), stype)
		case "proc":
			ss = steps.MakeCommandSource(cmd.MakeInteractiveCommand(spec_arg), stype)
		}
	default:
		Assert(false, "invalid -s, or -s not specified")
	}

	/*** -m ***/

	var m sample.Sampler
	spec_type, spec_flags, spec_arg = split(*mspec)
	switch spec_type {
	case "dummy":
		Assert(spec_flags == "" && spec_arg == "", "-m dummy does not accept flags/argument")
		m = sample.MakeDummySampler()
	case "cmd", "proc":
		Assert(spec_flags == "", "-m cmd/proc does not accept flags")
		switch spec_type {
		case "cmd":
			m = sample.MakeCommandSampler(cmd.MakeOneShotCommand(spec_arg))
		case "proc":
			m = sample.MakeCommandSampler(cmd.MakeInteractiveCommand(spec_arg))
		}
	default:
		Assert(false, "invalid -m, or -m not specified")
	}

	d := reg.MakeDomain(ts, ss, a, m)
	d.Start(fin, fout, throttle_type, throttle_period, t.Ticks(g))
	d.Wait()
}
