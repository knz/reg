package reg

import ("io"; "reg/ticks")

func MakeDomain(label string, ts ticks.Source) *Domain {
	dom := Domain {
		Label : label,
		ProtocolCmd : "while true; do read a || break; echo ACTION: $a >/dev/tty; done",
		TickSource : ts,

		resources : make(map[int]Resource),

		input : make(chan string),
		supplycmd : make(chan SupplyCmd),
		measure : make(chan Sample),
		query : make(chan bool),
		status : make(chan Status),
		action : make(chan Action),
		ticksctl : make(chan ticks.Ticks),
		statusctl : make(chan bool),
		tickssrc : make(chan ticks.Ticks),
		ticksext : make(chan ticks.Ticks),
		ticksin : make(chan ticks.Ticks),
		ticksper : make(chan ticks.Ticks),
		tickssteps : make(chan TicksSteps),
		stepsper : make(chan Steps),
		out : make(chan string),
		outready : make(chan bool),
		inputdone : make(chan bool) }


	return &dom
}

func (d *Domain) Start(input io.Reader) {
	d.TickSource.SetSource(d.ticksext)
	d.TickSource.Start()
	go d.readlines(input)
	go d.parse()
	go d.integrate()
	go d.protocol()
	go d.ticksource()
	go dupticks(d.ticksin, d.ticksper, d.tickssrc)
	go d.stepsource()
	go d.sample()
	go d.throttle()
	go d.outmgt()
	go d.output()
}

func (d *Domain) Wait() {
	<- d.inputdone
}


func (d *Domain) AddResource(label string, cmd string) {
	resnum := len(d.resources)
	d.resources[resnum] = Resource{label:label, cmd:cmd}
}
