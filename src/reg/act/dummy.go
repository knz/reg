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
