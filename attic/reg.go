package main

import (
	"strings"
	"strconv"
	"fmt"
	"bufio"
	"io"
	"os"
	"syscall"
	"time"
)

func report(outready <-chan bool, out chan<- string,
            supplyq chan<- bool, supplyr <-chan string,
            periodic <-chan bool, explicit <-chan bool) {
     for {
	     doit := false
	     select {
	     case <- outready:
	     case <- periodic: doit = true; continue
	     case <- explicit: doit = true; continue
	     }
	     if !doit {
		     select {
		     case <- periodic:
		     case <- explicit:
		     }
	     }

	     supplyq <- true;
	     response := <- supplyr;
	     out <- response;
	     doit = false
     }
}

type supplyorder struct {
        add bool;
	bin string;
	amount int;
};

func parse(cmd <-chan string, exp chan<- bool,
	ticks chan<- int, supply chan<- supplyorder) {
	for {
		cmdstr := <- cmd
		cmdargs := strings.Split(cmdstr, " ");

		switch cmdargs[0] {
		case ".":
			v, _ := strconv.ParseInt(cmdargs[1], 0, 64)
			ticks <- int(v)
		case "+":
			v := int64(0)
			if cmdargs[2] == "*" { v = -1
			} else { v, _ = strconv.ParseInt(cmdargs[2], 0, 64) }
			supply <- supplyorder{add:true, bin : cmdargs[1], amount : int(v)}
		case "-":
			v := int64(0)
			if cmdargs[2] == "*" {  v = -1
			} else { v, _ = strconv.ParseInt(cmdargs[2], 0, 64) }
			supply <- supplyorder{add:false, bin : cmdargs[1], amount : int(v)}
		case "?":
			exp <- true
		}
	}
}

func integrate(supply <-chan supplyorder, supplyq <-chan bool, supplyr chan<- string) {
	for {
		select {
		case order := <- supply:
			fmt.Println("Got an order %q", order)
		case <- supplyq:
			fmt.Println("Report request")
			supplyr <- "report response"
		}
	}
}


func readinput(file *bufio.Reader, cmd chan<- string, done chan<- bool) {
	for {
		line, err := file.ReadString('\n')
		if err == io.EOF { done <- true; break }
		cmd <- line
	}
}

func measure(ticks <-chan int) {
	for {
		select {
		case t := <- ticks: fmt.Println("Got ticks:", t)
		}
	}
}

func output(outready chan<- bool, out <-chan string) {
	set := syscall.FdSet{}
	for {
		set.Bits[0] = 2 // bit 1 = stdout
		syscall.Select(2, nil, &set, nil, nil);
		outready <- true;
		cmd := <- out;
		syscall.Write(1, []byte("Output " + cmd))
	}
}

func periodic(per chan<- bool) {
	for {
		per <- true
		time.Sleep(10 * time.Second)
	}
}

func main() {
	cmd := make(chan string)
	ri_done := make(chan bool)
	supply := make(chan supplyorder)
	ticks := make(chan int)
	supplyq := make(chan bool)
	supplyr := make(chan string)
	exp := make(chan bool)
	outready := make(chan bool)
	out := make(chan string)
	per := make(chan bool)
	go readinput(bufio.NewReader(os.Stdin), cmd, ri_done)
	go parse(cmd, exp, ticks, supply)
	go measure(ticks)
	go integrate(supply, supplyq, supplyr)
	go output(outready, out)
	go periodic(per)
	go report(outready, out, supplyq, supplyr, per, exp)



	<- ri_done
	fmt.Println("Done!")
}
