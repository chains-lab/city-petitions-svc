package petition

import (
	"context"

	svc "github.com/chains-lab/city-petitions-proto/gen/go/svc/petition"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/meta"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/problems"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/responses"
	"github.com/chains-lab/city-petitions-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) RejectPetition(ctx context.Context, req *svc.RejectPetitionRequest) (*svc.Petition, error) {
	initiator := meta.User(ctx)

	//TODO go to city svc check if user in city gov

	petitionId, err := uuid.Parse(req.GetPetitionId())
	if err != nil {
		logger.Log(ctx).Errorf("failed to parse petition id: %v", err)

		return nil, problems.InvalidArgumentError(ctx, "petition_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "petition_id",
			Description: "invalid UUID format for petition ID",
		})
	}

	petition, err := s.app.RejectPetition(ctx, petitionId, req.Reply)
	if err != nil {
		logger.Log(ctx).Errorf("failed to reject petition: %v", err)

		return nil, err
	}

	logger.Log(ctx).Infof("initiator %s is reject petition %s", initiator.ID, petitionId)

	return responses.Petition(petition), nil
}
