package main

import (
	"log"
	
	"pubsub/internal/delivery/grpc/server"
)

func main() {
	err := server.RunGRPCServer()
	if err != nil {
		log.Fatal(err)
	}
}
