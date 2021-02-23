first:
	go run main.go --name default1 --client-url localhost:3001 --peer-url localhost:5001 --peers "localhost:5002,localhost:5003"

second:
	go run main.go --name default2 --client-url localhost:3002 --peer-url localhost:5002 --peers "localhost:5001,localhost:5003"

third:
	go run main.go --name default3 --client-url localhost:3003 --peer-url localhost:5003 --peers "localhost:5001,localhost:5002"