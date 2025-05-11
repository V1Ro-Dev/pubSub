package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "pubsub/internal/delivery/grpc/proto"
	ps "pubsub/internal/subpub"
)

type PubSubServer struct {
	pb.UnimplementedPubSubServer
	psUseCase ps.SubPub
}

func NewPubSubServer(psUseCase ps.SubPub) *PubSubServer {
	return &PubSubServer{
		psUseCase: psUseCase,
	}
}

func (ps *PubSubServer) Subscribe(in *pb.SubscribeRequest, stream grpc.ServerStreamingServer[pb.Event]) error {
	ctx := stream.Context()
	reqID := ctx.Value("requestID")

	log.Printf("[reqID=%v] got subscribe request: %v", reqID, in)

	sub, err := ps.psUseCase.Subscribe(in.Key, func(msg interface{}) {
		if err := stream.Send(&pb.Event{
			Data: msg.(string),
		}); err != nil {
			log.Printf("Stream send error: %v", err)
		}
	})
	if err != nil {
		return status.Errorf(codes.Internal, "subscription error: %v", err)
	}

	<-ctx.Done()

	sub.Unsubscribe()

	if errors.Is(ctx.Err(), context.Canceled) {
		return nil
	}

	log.Printf("got subscription: %v", sub)
	return status.Errorf(codes.DeadlineExceeded, "context deadline exceeded")
}

func (ps *PubSubServer) Publish(ctx context.Context, in *pb.PublishRequest) (*emptypb.Empty, error) {
	reqID := ctx.Value("requestID")
	log.Printf("[reqID=%v] got publish request: %v", reqID, in)

	err := ps.psUseCase.Publish(in.Key, in.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to publish event: %v", err))
	}

	log.Printf("successfully proccessed publish request: %v", in)
	return &emptypb.Empty{}, nil
}
