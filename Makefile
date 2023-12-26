build:
	mkdir -p bin
	go build -o bin/pears-dht main.go

run:
	go run main.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/v1/chord.proto
	