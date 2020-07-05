package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"math"
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
	Value string
}

type SetResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func appendLog(command, key, value string) string {
	return command + ": " + key + " " + value
}

const (
	PUT    = "PUT"
	DELETE = "DELETE"
	GET    = "GET"
)

func CallRemoteNode(nodesToSendRPC []string, key string, value string) chan *AppendResult {
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
				appendArg.Entries = []string{appendLog("SET", key, value)}
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

	go func() {
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
		switch r.Method {
		case PUT:
		case DELETE:
		case GET:
			break
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("unsupported method"))
			return
		}
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

		coords.Log = append(coords.Log, appendLog("SET", updateKey.Key, updateKey.Value))

		appendResultChan := CallRemoteNode(nodesToSendRPC, updateKey.Key, updateKey.Value)

		successfulAppend := 0
		for ar := range appendResultChan {
			if ar.Success == true {
				successfulAppend++
			}
			// calculate that successful send is more than half
			if math.Ceil(float64(successfulAppend)/2) > float64(len(nodesToSendRPC)/2) {
				// TODO: add the value to the data store
			}
		}

		response, err := json.Marshal(SetResponse{Key: updateKey.Key, Value: updateKey.Value})

		if err != nil {
			log.Fatal("unable to marshal response JSON")
		}

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
