package errx

import (
	"context"
	"fmt"

	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/meta"
	"github.com/chains-lab/city-petitions-svc/internal/constant"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorInternal = ape.Declare("INTERNAL_ERROR")

func RaiseInternal(ctx context.Context, cause error) error {
	st := status.New(codes.Internal, "internal server error")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInternal.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorInternal.Raise(cause, st)
}

var ErrorRoleIsNotApplicable = ape.Declare("ROLE_IS_NOT_APPLICABLE")

func RaiseRoleIsNotApplicable(ctx context.Context, cause error, userID uuid.UUID, role string) error {
	msg := fmt.Sprintf("role is not applicable: user=%s role=%s", userID, role)
	st := status.New(codes.PermissionDenied, msg)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorRoleIsNotApplicable.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{RequestId: meta.RequestID(ctx)},
	)
	return ErrorRoleIsNotApplicable.Raise(cause, st)
}
