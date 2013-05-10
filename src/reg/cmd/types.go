package cmd

type Cmd interface {
	Start(in <-chan []string, out chan<- string)
}
