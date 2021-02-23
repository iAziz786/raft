package cmd

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

// SetKey will set the value for the respective key
func SetKey(coords *Coords) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case PUT:
		case DELETE:
			break
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("unsupported method"))
			return
		}
		var nodesToSendRPC []string

		for _, node := range peers {
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

		command := "PUT"
		coords.Log = append(coords.Log, appendLog(coords.Term, command, updateKey.Key, updateKey.Value))

		appendResultChan := CallRemoteNode(coords, nodesToSendRPC, command, updateKey.Key, updateKey.Value)

		successfulAppend := 1
		for ar := range appendResultChan {
			if ar.Success == true {
				successfulAppend++
			}
			// calculate that successful send is more than half
			if math.Ceil(float64(successfulAppend)/2) > float64(len(nodesToSendRPC)/2) {
				// add the value to the data store
				replicatedStateMachine[updateKey.Key] = updateKey.Value
			}
		}

		response, err := json.Marshal(LogEntry{Key: updateKey.Key, Value: updateKey.Value})

		if err != nil {
			log.Fatal("unable to marshal response JSON")
		}

		w.Header().Set("content-type", "application/json")
		_, err = w.Write(response)
		if err != nil {
			log.Println("unable to write the response")
		}

	}
}
