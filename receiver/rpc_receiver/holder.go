package rpc_receiver

import (
	"fmt"
	"time"

	"github.com/iAziz786/raft/config"
	"github.com/iAziz786/raft/receiver"
	"github.com/iAziz786/raft/replicator"
	"github.com/iAziz786/raft/storage"
)

type ServerState string

var (
	timer = time.NewTicker(1 * time.Second)
)

const (
	FOLLOWER  ServerState = "FOLLOWER"
	CANDIDATE             = "CANDIDATE"
	LEADER                = "LEADER"
)

type Coords struct {
	State    ServerState
	Term     uint32
	VotedFor interface{} // interface because it can be string or nil
	Log      []receiver.LogEntry
	storage  storage.Store
}

func NewCoords(storage storage.Store) *Coords {
	return &Coords{
		State:    FOLLOWER,
		Term:     0,
		VotedFor: nil,
		storage:  storage,
	}
}

type VoteArgument struct {
	Term         uint32
	CandidateId  string
	LastLogIndex uint32
	LastLogTerm  uint32
}

type VoteResult struct {
	Term        uint32
	VoteGranted bool
}

func (c *Coords) RequestVote(voteArg *VoteArgument, result *VoteResult) error {
	// 1. Reply false if term < currentTerm
	// 2. If votedFor is null or candidateId, and candidate's log is
	// at least as up-to-date as receiver's log, grand vote
	// check candidates log whether is it as up-to-date as mine
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
	return int(lastLogEntry.Term)
}

func (c *Coords) AppendEntry(appendArg *receiver.AppendArgument, state *replicator.AppendResult) error {
	timer.Reset(1 * time.Second)
	for _, entry := range appendArg.Entries {
		if err := c.storage.Set(entry.Key, entry.Value); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	fmt.Println("init holder")
	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("timer called")
				config.GetNotifier() <- struct{}{}
				fmt.Println("sent notification")
			default:
			}
		}
	}()
}
