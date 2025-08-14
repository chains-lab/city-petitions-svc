package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/problems"
	"github.com/chains-lab/city-petitions-svc/internal/constant"
	"github.com/chains-lab/city-petitions-svc/internal/logger"
	"github.com/chains-lab/gatekit/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ServiceJwtAuth(skService string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		switch info.FullMethod {
		case "add methods here if need":
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Log(ctx).Errorf("no metadata found in incoming context")

			return nil, problems.UnauthenticatedError(ctx, "no metadata found in incoming context")
		}

		token := md["x-service-token"]
		if len(token) == 0 {
			logger.Log(ctx).Errorf("service token not supplied")

			return nil, problems.UnauthenticatedError(ctx, fmt.Sprintf("service token not supplied"))
		}

		data, err := auth.VerifyServiceJWT(ctx, token[0], skService)
		if err != nil {
			logger.Log(ctx).Errorf("failed to verify service token: %s", err)

			return nil, problems.UnauthenticatedError(ctx, "failed to verify service token")
		}

		ThisSvcInAudience := false

		for _, aud := range data.Audience {
			if aud == constant.ServiceName {
				ThisSvcInAudience = true
				break
			}
		}

		if !ThisSvcInAudience {
			logger.Log(ctx).Errorf("service issuer %s not in audience %v", data.Issuer, data.Audience)

			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("service issuer %s not in audience %v", data.Issuer, data.Audience)).Err()
		}

		return handler(ctx, req)
	}
}
