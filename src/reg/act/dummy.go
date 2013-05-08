package act

type actuator_dummy struct{ actuator_common }

func MakeDummyActuator() Actuator {
	return &actuator_dummy{}
}

func (act *actuator_dummy) Start() {
	act.Check()

	go func() {
		for {
			// drain actions, do nothing
			<-act.source
		}
	}()
}
