package cmd

import "fmt"

type Coords struct {
	Term int
	Log  []int
}

func (c *Coords) Elect(name string, state *int) error {
	fmt.Println("calling elect with value", *state)

	return nil
}

func (c *Coords) AppendEntry(val int, state *int) error {
	fmt.Println("append entry")

	return nil
}