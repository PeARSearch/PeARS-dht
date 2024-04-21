build:
	mkdir -p bin
	go build -o bin/pears-dht main.go

run:
	go run main.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/v1/chord.proto

run-node-1:
	serf agent -config-file=./scripts/serf-config.json -event-handler="query:search_pears=./scripts/pears_query.sh" -node=node1 -bind=0.0.0.0:7946


run-node-2:
	serf agent -config-file=./scripts/serf-config.json -join=0.0.0.0:7946 -event-handler="query:search_pears=./scripts/pears_query-2.sh" -node=node2 -bind=0.0.0.0:7950 -rpc-addr=0.0.0.0:7374
