package petition

import (
	"context"

	svc "github.com/chains-lab/city-petitions-proto/gen/go/svc/petition"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/meta"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/problems"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/responses"
	"github.com/chains-lab/city-petitions-svc/internal/app/entities"
	"github.com/chains-lab/city-petitions-svc/internal/logger"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) CreatePetition(ctx context.Context, req *svc.CreatePetitionRequest) (*svc.Petition, error) {
	initiator := meta.User(ctx)

	cityID, err := uuid.Parse(req.CityId)
	if err != nil {
		return nil, problems.InvalidArgumentError(ctx, "city_id is invalid", &errdetails.BadRequest_FieldViolation{
			Field:       "city_id",
			Description: "invalid UUID format for city ID",
		})
	}

	petition, err := s.app.CreatePetition(ctx, cityID, initiator.ID, entities.CreatePetitionInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		logger.Log(ctx).Errorf("failed to create petition: %v", err)

		return nil, err
	}

	logger.Log(ctx).Infof("user %s created petition %s in city %s", initiator.ID, petition.ID, cityID)

	return responses.Petition(petition), nil
}
