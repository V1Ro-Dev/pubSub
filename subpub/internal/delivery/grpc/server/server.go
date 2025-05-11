package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"pubsub/internal/delivery/grpc/interceptors"
	pb "pubsub/internal/delivery/grpc/proto"
	ps "pubsub/internal/subpub"
)

func RunGRPCServer() error {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.ReqIDInterceptor),
		grpc.StreamInterceptor(interceptors.StreamReqIDInterceptor),
	)

	psUseCase := ps.NewPubSub()

	pb.RegisterPubSubServer(server, NewPubSubServer(&psUseCase))

	fmt.Println("starting server at :8081")
	err = server.Serve(lis)
	if err != nil {
		log.Fatalln(fmt.Sprintf("couldn't start server: %v", err))
		return err
	}

	return nil
}
