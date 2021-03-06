// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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
