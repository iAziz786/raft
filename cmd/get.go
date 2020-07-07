package cmd

import (
	"encoding/json"
	"log"
	"net/http"
)

type RequestValueWithKey struct {
	Key string
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body RequestValueWithKey
	decoder.Decode(&body)

	if val, ok := replicatedStateMachine[body.Key]; ok {
		response, err := json.Marshal(LogEntry{Key: body.Key, Value: val})

		if err != nil {
			log.Fatal("unable to marshal response JSON")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusNotFound)
		response, err := json.Marshal(LogEntry{Key: body.Key, Value: nil})

		if err != nil {
			log.Fatal("unable to marshal response JSON")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
