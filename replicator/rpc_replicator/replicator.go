package rpc_replicator

import (
	"log"
	"net/rpc"

	"github.com/iAziz786/raft/receiver"
	"github.com/iAziz786/raft/replicator"
)

type RPCReplicator struct {
	client map[string]*rpc.Client
}

func NewRPCReplicator() *RPCReplicator {
	return &RPCReplicator{
		client: map[string]*rpc.Client{},
	}
}

func (r RPCReplicator) Distribute(address, key string, msg []byte) error {
	client, err := r.getClient(address)
	if err != nil {
		return err
	}
	ar := replicator.AppendResult{}

	if err := client.Call("Coords.AppendEntry", receiver.AppendArgument{
		Term:         0,
		LeaderID:     "",
		PrevLogIndex: 0,
		PrevLogTerm:  0,
		Entries: []receiver.LogEntry{{
			Key:   key,
			Value: msg,
		}},
		LeaderCommitIndex: 0,
	}, &ar); err != nil {
		return err
	}

	return nil
}

func (r RPCReplicator) Elect(address string) error {
	log.Println("electing...")
	return nil
}

func (r RPCReplicator) getClient(address string) (*rpc.Client, error) {
	if client, ok := r.client[address]; ok {
		return client, nil
	}

	client, err := rpc.DialHTTP("tcp", address)

	if err != nil {
		return nil, err
	}

	r.client[address] = client
	return client, nil
}
