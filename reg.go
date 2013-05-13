package main

import . "assert"
import (
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
		Assert(err == nil, "-i :", err)
	}

	if *ofname == "-" {
		fout = os.Stdout
	} else {
		fout, err = os.OpenFile(*ofname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		Assert(err == nil, "-o :", err)
	}

	/*** -g ***/

	g, err := strconv.ParseFloat(*gran, 64)
	Assert(err == nil, "-g :", err)
	Assert(g >= 0, "-g : granularity cannot be negative")

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
		if spec_arg != "" {
			throttle_period, err = strconv.ParseFloat(spec_arg, 64)
			Assert(err == nil, "-p ticks :", err)
		}
		Assert(spec_flags == "", "-p ticks does not accept flags")
	case "steps":
		throttle_type = reg.OUTPUT_THROTTLE_STEPS
		if spec_arg != "" {
			throttle_period, err = strconv.ParseFloat(spec_arg, 64)
			Assert(err == nil, "-p steps :", err)
		}
		Assert(spec_flags == "", "-p steps does not accept flags")
	default:
		Assert(false, "invalid syntax for -p")
	}
	Assert(throttle_period >= 0, "-p : period cannot be negative")

	/*** -a ***/

	var a act.Actuator
	spec_type, spec_flags, spec_arg = split(*aspec)
	Assert(spec_flags == "", "-a does not accept flags")
	switch spec_type {
	case "discard":
		Assert(spec_arg == "", "-a discard does not accept argument")
		a = act.MakeDummyActuator()
	case "print":
		var af *os.File
		if spec_arg == "-" {
			af = os.Stdout
		} else {
			af, err = os.OpenFile(spec_arg, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			Assert(err == nil, "-a :", err)
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
	case "instant":
		Assert(spec_flags == "", "-t instant does not accept flags")
		v := float64(0)
		if spec_arg != "" {
			v, err = strconv.ParseFloat(spec_arg, 64)
			Assert(err == nil, "-t instant:", err)
		}
		ts = ticks.MakeDummySource(t.Ticks(v))
	case "time", "ptime", "cmd", "proc":
		stype, ok := flagconv(spec_flags)
		Assert(ok, "invalid flags for -t")
		switch spec_type {
		case "time", "ptime":
			d, err := time.ParseDuration(spec_arg)
			Assert(err == nil, "-t :", err)
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
	case "const":
		Assert(spec_flags == "", "-s dummy does not accept flags")
		v := float64(0)
		if spec_arg != "" {
			v, err = strconv.ParseFloat(spec_arg, 64)
			Assert(err == nil, "-s const:", err)
			Assert(v >= 0, "-s const: constant cannot be negative")
		}
		ss = steps.MakeDummySource(t.Steps(v))
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
	case "const":
		Assert(spec_flags == "", "-m const does not accept flags")
		v := float64(0)
		if spec_arg != "" {
			v, err = strconv.ParseFloat(spec_arg, 64)
			Assert(err == nil, "-m const:", err)
		}
		m = sample.MakeDummySampler(t.Stuff(v))
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
