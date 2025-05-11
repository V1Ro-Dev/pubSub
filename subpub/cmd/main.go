package main

import (
	"flag"
	"log"
	"pubsub/config"
	"pubsub/internal/delivery/grpc/server"
)

func main() {
	cfgPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg, err := config.Parse(*cfgPath)
	if err != nil {
		log.Printf("Config parse error occured, but default config was used: %v", err)
	}

	err = server.RunGRPCServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

}
