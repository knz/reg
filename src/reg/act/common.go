package act

import (
	"log"
	"reg/t"
)

type actuator_common struct {
	source chan t.Status
}

func (act *actuator_common) SetInput(src chan t.Status) {
	act.source = src
}

func (act *actuator_common) Check() {
	if act.source == nil {
		log.Fatal("no source channel connected")
	}
}
