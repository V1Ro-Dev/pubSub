package interceptors

import (
	"context"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ReqIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	const requestIDKey = "requestID"

	var reqID string
	if ids := md.Get(requestIDKey); len(ids) > 0 && ids[0] != "" {
		reqID = ids[0]
	} else {
		reqID = uuid.New().String()
		md.Set(requestIDKey, reqID)
	}

	ctx = metadata.NewIncomingContext(ctx, md)
	ctx = context.WithValue(ctx, requestIDKey, reqID)

	reply, err := handler(ctx, req)

	return reply, err
}

func StreamReqIDInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		md = metadata.New(nil)
	}

	const requestIDKey = "requestID"

	var reqID string
	if ids := md.Get(requestIDKey); len(ids) > 0 && ids[0] != "" {
		reqID = ids[0]
	} else {
		reqID = uuid.New().String()
		md.Set(requestIDKey, reqID)
	}

	// Добавляем в контекст
	ctx := context.WithValue(ss.Context(), requestIDKey, reqID)
	newStream := &wrappedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	}

	return handler(srv, newStream)
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
