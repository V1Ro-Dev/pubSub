package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"pubsub/config"
	"pubsub/internal/delivery/grpc/interceptors"
	pb "pubsub/internal/delivery/grpc/proto"
	ps "pubsub/internal/subpub"
)

func RunGRPCServer(cfg *config.Config) error {
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("can't listen on port: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.ReqIDInterceptor),
		grpc.StreamInterceptor(interceptors.StreamReqIDInterceptor),
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
		grpc.MaxConcurrentStreams(cfg.MaxConcurrentStreams),
	)

	psUseCase := ps.NewPubSub()
	pb.RegisterPubSubServer(grpcServer, NewPubSubServer(&psUseCase))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Printf("starting server at %s\n", cfg.Addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("couldn't start server: %v", err)
		}
	}()

	<-stop
	fmt.Println("\nshutting down server gracefully...")

	grpcServer.GracefulStop()
	fmt.Println("server stopped")

	return nil
}
