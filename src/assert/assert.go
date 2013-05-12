package assert

import (
	"fmt"
	"os"

//	"runtime/debug"
)

func CheckErrIsNil(err error, ctx ...interface{}) {
	if err != nil {
		// fmt.Fprintln(os.Stderr, "Stack trace:")
		// debug.PrintStack()
		fmt.Fprint(os.Stderr, ctx...)
		fmt.Fprintln(os.Stderr, ":", err)
		os.Exit(1)
	}
}
