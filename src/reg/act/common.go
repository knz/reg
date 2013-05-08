package act

import "log"

type actuator_common struct {
	source chan Action
}

func (act *actuator_common) SetInput(src chan Action) {
	act.source = src
}

func (act *actuator_common) Check() {
	if act.source == nil {
		log.Fatal("no source channel connected")
	}
}
