# raft
The Raft implementation in Go

### How to run the program?

All the nodes that are going to be  in the cluster should be mentioned in the startup command.
```shell script
go run main.go --http-port 3001 --raft-port 5001 --nodes "localhost:3001,localhost:3002,localhost:3003" 
```

The `http-port` is the port where the communication happen of the node gets elected as the leader. The `raft-port` is
used to make the RPC call. You need to provide all the node that would be in the cluster in the `node` flag.

### 🚧 Feature working in progress

- [ ] Leader election
- [ ] Membership change
- [ ] Log compaction