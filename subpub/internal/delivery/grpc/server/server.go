package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"pubsub/config"
	"pubsub/internal/delivery/grpc/interceptors"
	pb "pubsub/internal/delivery/grpc/proto"
	ps "pubsub/internal/subpub"
)

func RunGRPCServer(cfg *config.Config) error {
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.ReqIDInterceptor),
		grpc.StreamInterceptor(interceptors.StreamReqIDInterceptor),
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
		grpc.MaxConcurrentStreams(cfg.MaxConcurrentStreams),
	)

	psUseCase := ps.NewPubSub()

	pb.RegisterPubSubServer(server, NewPubSubServer(&psUseCase))

	fmt.Println(fmt.Sprintf("starting server at %s", cfg.Addr))
	err = server.Serve(lis)
	if err != nil {
		log.Fatalln(fmt.Sprintf("couldn't start server: %v", err))
		return err
	}

	return nil
}
