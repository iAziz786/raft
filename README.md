# raft
The Raft implementation in Go

### How to run the program?

All the nodes that are going to be  in the cluster should be mentioned in the startup command.
```shell script
go run main.go --name default1 --client-url localhost:3001 --peer-url localhost:5001 --peers "localhost:5002,localhost:5003"
```

The `client-url` is the port where the communication happen of the node gets elected as the leader. The `peer-url` is
used to make the RPC call. You need to provide all the node that would be in the cluster in the `node` flag.

### 🚧 Feature working in progress

- [ ] Leader election
- [ ] Membership change
- [ ] Log compaction