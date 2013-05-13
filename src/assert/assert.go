package assert

import (
	"fmt"
	"os"

//	"runtime/debug"
)

func Assert(cond bool, ctx ...interface{}) {
	if !cond {
		fmt.Fprintln(os.Stderr, ctx...)
		os.Exit(1)
	}
}

func CheckErrIsNil(err error, ctx ...interface{}) {
	if err != nil {
		// fmt.Fprintln(os.Stderr, "Stack trace:")
		// debug.PrintStack()
		fmt.Fprint(os.Stderr, ctx...)
		fmt.Fprintln(os.Stderr, ":", err)
		os.Exit(1)
	}
}
