package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "raft",
	Short: "raft uses distributed consensus algorithm",
	Long:  `Raft paper was created by one PhD researcher at stanford university`,
	Run:   Run,
}

var httpPort string
var raftPort string

// nodes are the list tpc address which serves RPC calls
var nodes []string

type UpdateKey struct {
	Key   string
	Value string
}

type LogEntry struct {
	Term    int         `json:"-"`
	Command string      `json:"-"`
	Key     string      `json:"key"`
	Value   interface{} `json:"value"`
}

var replicatedStateMachine = make(map[string]string)

func appendLog(term int, command, key, value string) LogEntry {
	return LogEntry{
		Term:    term,
		Key:     key,
		Value:   value,
		Command: command,
	}
}

const (
	PUT    = "PUT"
	DELETE = "DELETE"
	GET    = "GET"
)

func CallRemoteNode(coords *Coords, nodesToSendRPC []string, command, key, value string) chan *AppendResult {
	appendResult := make(chan *AppendResult)
	var wg sync.WaitGroup
	go func() {
		defer close(appendResult)
		for _, node := range nodesToSendRPC {
			wg.Add(1)
			go func(node string) {
				defer wg.Done()
				fmt.Println("dialing client", node)
				client, err := rpc.DialHTTP("tcp", node)
				if err != nil {
					log.Fatal("dialing error:", err)
				}

				var appendResultForThisNode AppendResult
				var appendArg AppendArgument

				appendArg.Term = 1
				appendArg.Entries = []LogEntry{appendLog(coords.Term, command, key, value)}
				appendArg.LeaderCommitIndex = 1
				appendArg.LeaderId = raftPort
				appendArg.PrevLogIndex = 1
				appendArg.PrevLogTerm = 1

				err = client.Call("Coords.AppendEntry", &appendArg, &appendResultForThisNode)
				if err != nil {
					log.Println("error while calling the elect", err)
				}

				appendResult <- &appendResultForThisNode
			}(node)
		}
		wg.Wait()
	}()
	return appendResult
}

func Run(cmd *cobra.Command, args []string) {
	coords := NewCoords()
	err := rpc.Register(coords)
	if err != nil {
		log.Fatalf("unable to register the struct")
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+raftPort)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	go func() {
		err = http.Serve(l, nil)
		if err != nil {
			log.Fatal("serving error:", err)
		}
	}()

	http.HandleFunc("/value", GetKey)

	http.HandleFunc("/set", SetKey(coords))

	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		fmt.Println("error while serving", err)
		os.Exit(1)
	}
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&httpPort, "http-port", "p", "", "run the http server to handle the clients")
	rootCmd.PersistentFlags().StringVarP(&raftPort, "raft-port", "r", "", "communicate with other rpc servers on this port")
	rootCmd.PersistentFlags().StringSliceVarP(&nodes, "nodes", "n", []string{}, "endpoint for all the nodes in the cluster")

	if rootCmd.MarkPersistentFlagRequired("http-port") != nil {
		log.Fatalf("unable to make flag %s required", "http-port")
	}
	if rootCmd.MarkPersistentFlagRequired("raft-port") != nil {
		log.Fatalf("unable to make flag %s required", "raft-port")
	}
	if rootCmd.MarkPersistentFlagRequired("nodes") != nil {
		log.Fatalf("unable to make flag %s required", "nodes")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
