package cmd

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sync"
	"time"

	"github.com/iAziz786/raft/config"
	"github.com/iAziz786/raft/receiver/rpc_receiver"
	"github.com/iAziz786/raft/replicator/rpc_replicator"
	"github.com/iAziz786/raft/server"
	"github.com/iAziz786/raft/storage/bbolt_store"
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

// peers are the list tpc address which serves RPC calls
var peers []string

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
	nrr := rpc_receiver.NewRPCReceiver()
	rpcRepl := rpc_replicator.NewRPCReplicator()
	go nrr.Receive()
	go func() {
		for {
			select {
			case <-config.GetNotifier():
				fmt.Println("notified")
				for _, peer := range config.Peers {
					rpcRepl.Elect(peer)
				}
			default:
				time.Sleep(100 * time.Millisecond)
				fmt.Println("default")
			}
		}
	}()
	server := server.NewServer(rpcRepl, bbolt_store.NewStore())
	server.Serve(config.ClientURL)
}

// Execute validates and execute the commands
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&config.ClientURL, "client-url", "p", "", "run the http server to handle the clients")
	rootCmd.PersistentFlags().StringVarP(&config.PeerURL, "peer-url", "r", "", "communicate with other rpc servers on this port")
	rootCmd.PersistentFlags().StringSliceVarP(&config.Peers, "peers", "n", []string{}, "endpoint for all the nodes in the cluster")
	rootCmd.PersistentFlags().StringVarP(&config.Name, "name", "i", "default", "name of the node to identify it")

	if rootCmd.MarkPersistentFlagRequired("client-url") != nil {
		log.Fatalf("unable to make flag %s required", "client-url")
	}
	if rootCmd.MarkPersistentFlagRequired("peer-url") != nil {
		log.Fatalf("unable to make flag %s required", "peer-url")
	}
	if rootCmd.MarkPersistentFlagRequired("peers") != nil {
		log.Fatalf("unable to make flag %s required", "peers")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
