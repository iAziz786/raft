package cmd

import "fmt"

type ServerState string

const (
	FOLLOWER  ServerState = "FOLLOWER"
	CANDIDATE             = "CANDIDATE"
	LEADER                = "LEADER"
)

type Coords struct {
	State    ServerState
	Term     int
	VotedFor interface{} // interface because it can be string or nil
	Log      []LogEntry
}

func NewCoords() *Coords {
	return &Coords{
		State:    FOLLOWER,
		Term:     0,
		VotedFor: nil,
	}
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
	Entries           []LogEntry
	LeaderCommitIndex int
}

type VoteArgument struct {
	Term         int
	CandidateId  string
	LastLogIndex int
	LastLogTerm  int
}

type VoteResult struct {
	Term        int
	VoteGranted bool
}

func (c *Coords) RequestVote(voteArg *VoteArgument, result *VoteResult) error {
	// 1. Reply false if term < currentTerm
	if voteArg.Term < c.Term {
		result.Term = c.Term
		result.VoteGranted = false
		return nil
	}
	// 2. If votedFor is null or candidateId, and candidate's log is
	// at least as up-to-date as receiver's log, grand vote
	if c.VotedFor == nil || c.VotedFor == voteArg.CandidateId {
		// check candidates log whether is it as up-to-date as mine
		if c.getLastLogIndex() <= voteArg.LastLogIndex && c.getLastLogTerm() <= voteArg.LastLogTerm {
			result.Term = voteArg.Term
			result.VoteGranted = true
			return nil
		}
	}
	result.VoteGranted = false
	return nil
}

func (c *Coords) getLastLogIndex() int {
	if len(c.Log) == 0 {
		return 0
	}
	return len(c.Log) - 1
}

func (c *Coords) getLastLogTerm() int {
	if len(c.Log) == 0 {
		return 0
	}

	lastLogEntry := c.Log[len(c.Log)-1]
	return lastLogEntry.Term
}

func (c *Coords) AppendEntry(appendArg *AppendArgument, state *AppendResult) error {
	fmt.Println("append entry")
	c.Log = append(c.Log, appendArg.Entries...)

	fmt.Println(c.Log)

	state.Success = true
	state.Term = c.Term

	return nil
}
