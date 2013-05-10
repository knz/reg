package act

import "reg/t"

type Actuator interface {
	Start(src <-chan t.Status)
}
