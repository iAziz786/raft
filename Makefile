first:
	go run main.go --http-port 3001 --raft-port 5001 --nodes "localhost:3001,localhost:3002,localhost:3003"

second:
	go run main.go --http-port 3002 --raft-port 5002 --nodes "localhost:3001,localhost:3002,localhost:3003"

third:
	go run main.go --http-port 3003 --raft-port 5003 --nodes "localhost:3001,localhost:3002,localhost:3003"