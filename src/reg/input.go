// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package reg

import . "assert"
import (
	"bufio"
	"io"
	"reg/t"
	"strconv"
	"strings"
)

func readlines(input io.Reader, dst chan<- string, inputdone chan<- bool) {
	reader := bufio.NewReader(input)
	for {
		cmdstr, err := reader.ReadString('\n')
		if err == io.EOF {
			inputdone <- true
			break
		}
		dst <- cmdstr[:len(cmdstr)-1]
	}

}

func parse(input <-chan string, ticksctl chan<- t.Ticks, supplycmd chan<- SupplyCmd, acmd chan<- bool, statusctl chan<- bool) {
	for cmd := range input {
		cmdargs := strings.Split(cmd, " ")

		switch cmdargs[0] {
		case ".":
			Assert(len(cmdargs) == 2, "invalid syntax for . on input: ", cmd)
			v, err := strconv.ParseFloat(cmdargs[1], 64)
			Assert(err == nil, "parsing . on input", ":", err)
			ticksctl <- t.Ticks(v)
		case "+":
			Assert(len(cmdargs) == 2, "invalid syntax for + on input: ", cmd)
			v, err := strconv.ParseFloat(cmdargs[1], 64)
			Assert(err == nil, "parsing + on input", ":", err)
			supplycmd <- SupplyCmd{t.StuffSteps(v)}
		case "aon":
			acmd <- true
		case "aoff":
			acmd <- false
		case "?":
			Assert(len(cmdargs) == 1, "invalid syntax for ? on input: ", cmd)
			statusctl <- true
		}
	}
}
