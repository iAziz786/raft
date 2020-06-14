package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "raft",
	Short: "raft uses distributed consensus algorithm",
	Long:  `Raft paper was created by one PhD researcher at stanford university`,
	Run:   Run,
}

var httpPort string
var raftPort string
var nodes []string

func Run(cmd *cobra.Command, args []string) {
	go func() {
		coords := new(Coords)
		rpc.Register(coords)
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
