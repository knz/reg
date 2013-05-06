package reg

type Ticks float64
type Steps float64
type Stuff float64
type StuffSteps float64

type TicksSteps struct {
	ticks Ticks
	steps Steps
}

type SupplyCmd struct {
	bin int;
	supply StuffSteps;
}

type Sample struct {
        ticks Ticks;
	steps Steps;
	usage []Stuff;
}

type Status struct {
	ticks Ticks;
	steps Steps;
	usage []StuffSteps;
}

type Action struct {
	bin int;
	currentsupply StuffSteps;
	delta StuffSteps;
}

type Resource struct {
	label string
	cmd string
}

const ( ThrottleSteps = iota; ThrottleTicks )

type TickSource interface {
	Start()
	GetSource() chan Ticks
}

type Domain struct {
	Label string
	TickSource TickSource
	StepsCmd string
	ProtocolCmd string
	OutputFile string

	ThrottleType int
	ThrottleMinPeriod float64

	// Resource management
	resources map[int]Resource

	// Channels
	input chan string // readlines -> parse
	supplycmd chan SupplyCmd // parse -> integrate
	measure chan Sample // sample -> integrate
	query chan bool // outputmgt -> integrate
	status chan Status // integrate -> outputmgt
	action chan Action // integrate -> protocol
	ticksctl chan Ticks // parse -> ticksource
	statusctl chan bool // parse -> outmgt
	tickssrc chan Ticks // ticksource -> dup
	ticksin chan Ticks // dup -> stepsource
	ticksper chan Ticks // dup -> throttle
	tickssteps chan TicksSteps // stepsource -> sample
	stepsper chan Steps // stepsource -> throttle

	out chan string // outmgt -> output
	outready chan bool // output -> outmgt

	inputdone chan bool // parse -> wait
}
