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

var ErrorPetitionSignaturesNotFound = ape.Declare("PETITION_SIGNATURES_NOT_FOUND")

func RaisePetitionSignaturesNotFoundByPetitionIDUserID(ctx context.Context, cause error, petitionID, userID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("Petition signatures for petition ID '%s' and user ID '%s' not found", petitionID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorPetitionSignaturesNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx), // Request ID is not available in this context
		},
	)

	return ErrorPetitionSignaturesNotFound.Raise(cause, st)
}

func RaisePetitionSignaturesNotFoundByID(ctx context.Context, cause error, sigID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("Petition signatures with ID '%s' not found", sigID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorPetitionSignaturesNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx), // Request ID is not available in this context
		},
	)

	return ErrorPetitionSignaturesNotFound.Raise(cause, st)
}

var ErrorPetitionSignaturesAlreadyExists = ape.Declare("PETITION_SIGNATURES_ALREADY_EXISTS")

func RaisePetitionSignaturesAlreadyExists(ctx context.Context, cause error, petitionID, userID uuid.UUID) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("Petition signatures for petition ID '%s' and user ID '%s' already exists", petitionID, userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorPetitionSignaturesAlreadyExists.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx), // Request ID is not available in this context
		},
	)

	return ErrorPetitionSignaturesAlreadyExists.Raise(cause, st)
}
