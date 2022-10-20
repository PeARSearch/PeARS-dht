FROM golang:1.19.2-alpine3.16 AS build

RUN mkdir /pears
WORKDIR /pears

COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY pkg pkg

RUN go build -o /pears-dht ./cmd/pears-dht/main.go

FROM alpine:3.16 as deploy

ENV CDHT_PORT=8080
ENV CDHT_SERVER_PORT=8888
ENV CDHT_CONTACTS=
ENV CDHT_SEED=14348

COPY --from=build /pears-dht /dht

RUN addgroup PeARS && \
    adduser -u 1001 -D -G PeARS cloud && \
    chown cloud:PeARS /dht

USER 1001

ENTRYPOINT ["./dht"]
