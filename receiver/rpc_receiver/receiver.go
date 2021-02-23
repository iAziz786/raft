package rpc_receiver

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/iAziz786/raft/config"
	"github.com/iAziz786/raft/storage/bbolt_store"
)

type RPCReceiver struct {
	coords *Coords
}

func NewRPCReceiver() *RPCReceiver {
	return &RPCReceiver{
		coords: &Coords{},
	}
}

func (r RPCReceiver) Receive() error {
	if err := rpc.Register(NewCoords(bbolt_store.NewStore())); err != nil {
		log.Println("unable to register RPC")
		return err
	}

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", config.PeerURL)
	if err != nil {
		log.Printf("listening error: %s\n", err)
		return err
	}

	return http.Serve(l, nil)
}
