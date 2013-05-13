// Copyright 2013 Raphael 'kena' Poss.  All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package act

import "reg/t"

type actuator_dummy struct{}

func MakeDummyActuator() Actuator {
	return &actuator_dummy{}
}

func (act *actuator_dummy) Start(src <-chan t.Status) {
	for {
		// drain actions, do nothing
		<-src
	}
}
