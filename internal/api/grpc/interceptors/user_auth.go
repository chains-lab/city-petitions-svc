package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/meta"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/problems"
	"github.com/chains-lab/city-petitions-svc/internal/logger"
	"github.com/chains-lab/gatekit/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UserJwtAuth(skUser string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Log(ctx).Errorf("no metadata found in incoming context")

			return nil, problems.UnauthenticatedError(ctx, "no metadata found in incoming context")
		}

		token := md["x-user-token"]
		if len(token) == 0 {
			logger.Log(ctx).Errorf("user token not supplied")

			return nil, problems.UnauthenticatedError(ctx, fmt.Sprintf("user token not supplied"))
		}

		userData, err := auth.VerifyUserJWT(ctx, token[0], skUser)
		if err != nil {
			logger.Log(ctx).Errorf("failed to verify user token: %s", err)

			return nil, problems.UnauthenticatedError(ctx, "failed to verify user token")
		}

		userID, err := uuid.Parse(userData.Subject)
		if err != nil {
			logger.Log(ctx).Errorf("invalid user ID: %v", err)

			return nil, problems.UnauthenticatedError(ctx, fmt.Sprintf("invalid user ID: %v", err))
		}

		ctx = context.WithValue(ctx, meta.UserCtxKey, meta.UserData{
			ID:        userID,
			SessionID: userData.Session,
			Verified:  userData.Verified,
			Role:      userData.Role,
		})

		return handler(ctx, req)
	}
}
