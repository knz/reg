package t

type Ticks float64
type Steps float64
type Stuff float64
type StuffSteps float64

const (
	SRC_MONOTONIC = iota
	SRC_INIT_THEN_DELTAS
	SRC_DELTAS_ONLY
)

type TicksSteps struct {
	Ticks Ticks
	Steps Steps
}

type Status struct {
	Ticks      Ticks
	TicksDelta Ticks
	Steps      Steps
	StepsDelta Steps
	Supply     StuffSteps
	Delta      StuffSteps
}

type Sample struct {
	Ticks Ticks
	Steps Steps
	Usage Stuff
}
