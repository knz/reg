package assert

import (
	"fmt"
	"os"
)

func Assert(cond bool, ctx ...interface{}) {
	if !cond {
		fmt.Fprintln(os.Stderr, ctx...)
		os.Exit(1)
	}
}
