package errx

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/meta"
	"github.com/chains-lab/city-petitions-svc/internal/constant"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func nowRFC3339Nano() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

var ErrorPetitionNotFound = ape.Declare("PETITION_NOT_FOUND")

func RaisePetitionNotFoundByID(ctx context.Context, cause error, petitionID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("Petition with id '%s' not found", petitionID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorPetitionNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx), // Request ID is not available in this context
		},
	)

	return ErrorPetitionNotFound.Raise(cause, st)
}

var ErrorPetitionIsNotAvailable = ape.Declare("PETITION_IS_NOT_AVAILABLE")

func RaisePetitionIsNotAvailable(ctx context.Context, cause error, petitionID string) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("Petition with id '%s' is not available", petitionID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorPetitionIsNotAvailable.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
		&errdetails.RequestInfo{
			RequestId: meta.RequestID(ctx), // Request ID is not available in this context
		},
	)

	return ErrorPetitionIsNotAvailable.Raise(cause, st)
}
