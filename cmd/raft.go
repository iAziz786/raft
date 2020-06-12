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

func Run(cmd *cobra.Command, args []string) {
	go func() {
		coords := new(Coords)
		rpc.Register(coords)
		rpc.HandleHTTP()
		l, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatal("listen error:", err)
		}

		err = http.Serve(l, nil)
		if err != nil {
			log.Fatal("serving error:", err)
		}
	}()

	http.HandleFunc("/key", func(writer http.ResponseWriter, r *http.Request) {
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Fatal("dialing error:", err)
		}

		client.Call("Coords.Elect", "rambo", 4)
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("error while serving", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
