package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/meta"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/problems"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, problems.UnauthenticatedError(ctx, fmt.Sprintf("no metadata found in incoming context"))
		}

		requestIDArr := md["x-request-id"]
		if len(requestIDArr) == 0 {
			return nil, problems.UnauthenticatedError(ctx, fmt.Sprintf("request ID not supplied"))
		}

		requestID, err := uuid.Parse(requestIDArr[0])
		if err != nil {
			return nil, problems.UnauthenticatedError(ctx, fmt.Sprintf("invalid request ID: %v", err))
		}

		ctx = context.WithValue(ctx, meta.RequestIDCtxKey, requestID)

		return handler(ctx, req)
	}
}
