package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iAziz786/raft/config"
	"github.com/iAziz786/raft/replicator"
	"github.com/iAziz786/raft/storage"
)

type Server struct {
	replicator replicator.Replicator
	storage    storage.Store
}

type setKeyReq struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func NewServer(replicator replicator.Replicator, storage storage.Store) *Server {
	return &Server{
		replicator: replicator,
		storage:    storage,
	}
}

func (s *Server) Serve(address string) {
	r := mux.NewRouter()
	r.HandleFunc("/store", func(w http.ResponseWriter, req *http.Request) {
		var skReq setKeyReq
		d := json.NewDecoder(req.Body)
		err := d.Decode(&skReq)
		if err != nil {
			panic(err)
		}

		for _, address := range config.Peers {
			if err := s.replicator.Distribute(address, skReq.Key, []byte(skReq.Val)); err != nil {
				log.Println(err)
			}
		}

		if err := s.replicator.Distribute(config.PeerURL, skReq.Key, []byte(skReq.Val)); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			log.Println(err)
			w.Write([]byte("failed to distribute the message"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("POST")

	r.HandleFunc("/store/{key}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		key := vars["key"]
		val, err := s.storage.Get(key)
		if err != nil {
			if err == storage.ErrNotFound {
				w.WriteHeader(http.StatusOK)
				w.Write(nil)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(val)
	}).Methods("GET")
	http.ListenAndServe(address, r)
}
