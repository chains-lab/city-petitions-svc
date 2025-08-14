package errx

import (
	"context"
	"fmt"

	"github.com/chains-lab/city-petitions-svc/internal/api/meta"
	"github.com/chains-lab/city-petitions-svc/internal/constant"
	"github.com/chains-lab/svc-errors/ape"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorPetitionSignaturesNotFound = ape.Declare("PETITION_SIGNATURES_NOT_FOUND")

func RaisePetitionSignaturesNotFoundByID(ctx context.Context, cause error, petitionID string) error {
	st := status.New(codes.NotFound, fmt.Sprintf("petition signatures with petition ID '%s' not found", petitionID))
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
