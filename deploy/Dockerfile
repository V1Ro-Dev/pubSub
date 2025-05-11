FROM golang:1.24 AS builder

WORKDIR /subpub_app

COPY subpub/go.mod subpub/go.sum subpub/
RUN cd subpub && go mod download

COPY subpub /subpub_app/subpub

WORKDIR /subpub_app/subpub/cmd

RUN go build -o main main.go

FROM debian:bookworm-slim

WORKDIR /subpub_app/

COPY --from=builder /subpub_app/subpub/cmd/main .
COPY deploy/config/server.yml ./deploy/config/server.yml

CMD ["./main", "-config=./deploy/config/server.yml"]
