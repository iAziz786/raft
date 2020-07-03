package cmd

import "fmt"

type Coords struct {
	Term int
	Log  []int
}

type AppendResult struct {
	Term    int
	Success bool
}

type AppendArgument struct {
	Term              int
	LeaderId          string
	PrevLogIndex      int
	PrevLogTerm       int
	Entries           []string
	LeaderCommitIndex int
}

func (c *Coords) Elect(appendArg *AppendArgument, state *AppendResult) error {
	fmt.Printf("electing from %s\n", httpPort)

	return nil
}

func (c *Coords) AppendEntry(appendArg *AppendArgument, state *AppendResult) error {
	fmt.Println("append entry")
	fmt.Println(appendArg)

	state.Success = true
	state.Term = c.Term

	return nil
}
