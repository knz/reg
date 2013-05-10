package reg

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

func parse(input <-chan string, ticksctl chan<- t.Ticks, supplycmd chan<- SupplyCmd, statusctl chan<- bool) {
	for cmd := range input {
		cmdargs := strings.Split(cmd, " ")

		switch cmdargs[0] {
		case ".":
			v, _ := strconv.ParseFloat(cmdargs[1], 64)
			ticksctl <- t.Ticks(v)
		case "+":
			b, _ := strconv.ParseInt(cmdargs[1], 0, 0)
			v, _ := strconv.ParseFloat(cmdargs[2], 64)
			supplycmd <- SupplyCmd{bin: int(b), supply: t.StuffSteps(v)}
		case "?":
			statusctl <- true
		}
	}
}
