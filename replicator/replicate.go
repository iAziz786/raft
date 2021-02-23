package replicator

type AppendResult struct {
	Term    uint32
	Success bool
}

type Replicator interface {
	Distribute(address, key string, msg []byte) error
}
