package act

import "reg/t"

type Actuator interface {
	Start()
	SetInput(src chan t.Status)
}
