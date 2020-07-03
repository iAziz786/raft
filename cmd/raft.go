package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
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
	Value interface{}
}

type SetResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func CallRemoteNode(nodesToSendRPC []string) chan *AppendResult {
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
				appendArg.Entries = []string{"set value 1"}
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
	go func() {
		coords := new(Coords)
		err := rpc.Register(coords)
		if err != nil {
			log.Fatalf("unable to register the struct")
		}
		rpc.HandleHTTP()
		l, err := net.Listen("tcp", ":"+raftPort)
		if err != nil {
			log.Fatal("listen error:", err)
		}

		err = http.Serve(l, nil)
		if err != nil {
			log.Fatal("serving error:", err)
		}
	}()

	http.HandleFunc("/key", func(w http.ResponseWriter, r *http.Request) {
		client, err := rpc.DialHTTP("tcp", pickRandomElement(nodes))
		if err != nil {
			log.Fatal("dialing error:", err)
		}

		client.Call("Coords.Elect", "rambo", 4)
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		var nodesToSendRPC []string

		for _, node := range nodes {
			if "localhost:"+raftPort != node {
				nodesToSendRPC = append(nodesToSendRPC, node)
			}
		}

		decoder := json.NewDecoder(r.Body)

		var updateKey UpdateKey

		err := decoder.Decode(&updateKey)

		if err != nil {
			log.Fatalf("failed to decode updateKey")
		}

		appendResultChan := CallRemoteNode(nodesToSendRPC)

		for ar := range appendResultChan {
			fmt.Println("term", ar.Term)
			fmt.Println("success", ar.Success)
		}

		response, err := json.Marshal(SetResponse{Key: "Yolo", Value: true})

		w.Header().Set("content-type", "application/json")
		_, err = w.Write(response)
		if err != nil {
			log.Println("unable to write the response")
		}
	})

	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		fmt.Println("error while serving", err)
		os.Exit(1)
	}
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&httpPort, "http-port", "p", "", "run the http server to handle the clients")
	rootCmd.PersistentFlags().StringVarP(&raftPort, "raft-port", "r", "", "communicate with other rpc servers on this port")
	rootCmd.PersistentFlags().StringSliceVarP(&nodes, "nodes", "n", []string{}, "endpoint for all the nodes in the cluster")

	rootCmd.MarkPersistentFlagRequired("http-port")
	rootCmd.MarkPersistentFlagRequired("raft-port")
	rootCmd.MarkPersistentFlagRequired("nodes")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
