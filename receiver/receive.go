package receiver

type LogEntry struct {
	Term    uint32
	Command string
	Key     string
	Value   []byte
}

type AppendArgument struct {
	Term              uint32
	LeaderID          string
	PrevLogIndex      uint32
	PrevLogTerm       uint32
	Entries           []LogEntry
	LeaderCommitIndex uint32
}

type Receiver interface {
	Receive(key string, msg []byte) error
}
