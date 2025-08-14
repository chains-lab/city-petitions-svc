package petition

import (
	"context"

	svc "github.com/chains-lab/city-petitions-proto/gen/go/svc/petition"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/problems"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/responses"
	"github.com/chains-lab/city-petitions-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetPetition(ctx context.Context, req *svc.GetPetitionRequest) (*svc.Petition, error) {
	petitionID, err := uuid.Parse(req.GetPetitionId())
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "petition_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "petition_id",
			Description: "invalid UUID format for petition ID",
		})
	}

	petition, err := s.app.GetPetition(ctx, petitionID)
	if err != nil {
		logger.Log(ctx).Errorf("failed to get petition: %v", err)

		return nil, err
	}

	return responses.Petition(petition), nil
}
