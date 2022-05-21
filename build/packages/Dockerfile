FROM golang:alpine3.15

RUN mkdir /pears
WORKDIR /pears

COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY pkg pkg

RUN go build -o pears-dht ./cmd/cacophany-dht/main.go

ENTRYPOINT ["/pears/pears-dht -p 8080"]
