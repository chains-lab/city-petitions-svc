package petition

import (
	"context"

	svc "github.com/chains-lab/city-petitions-proto/gen/go/svc/petition"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/problems"
	"github.com/chains-lab/city-petitions-svc/internal/api/grpc/responses"
	"github.com/chains-lab/city-petitions-svc/internal/app/entities"
	"github.com/chains-lab/city-petitions-svc/internal/logger"
	"github.com/chains-lab/city-petitions-svc/internal/pagination"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) ListPetitionSigners(ctx context.Context, req *svc.ListPetitionSignersRequest) (*svc.SignatureList, error) {
	filters := entities.ListPetitionsSignFilter{}

	if req.PetitionId != nil {
		petitionID, err := uuid.Parse(*req.PetitionId)
		if err != nil {
			logger.Log(ctx).Errorf("failed to parse petition id: %v", err)
			return nil, problems.InvalidArgumentError(ctx, "petition_id is invalid", &errdetails.BadRequest_FieldViolation{
				Field:       "petition_id",
				Description: "invalid UUID format for petition ID",
			})
		}

		filters.PetitionID = &petitionID
	}

	if req.UserId != nil {
		userID, err := uuid.Parse(*req.UserId)
		if err != nil {
			logger.Log(ctx).Errorf("failed to parse user id: %v", err)
			return nil, problems.InvalidArgumentError(ctx, "user_id is invalid", &errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "invalid UUID format for user ID",
			})
		}

		filters.UserID = &userID
	}

	sort := entities.ListPetitionsSignSort{}
	switch srt := req.Sort.(type) {
	case *svc.ListPetitionSignersRequest_Oldest:
		if srt.Oldest {
			sort.Oldest = true
		} else {
			sort.Oldest = true
		}
	default:
		sort.Newest = true
	}

	signers, pag, err := s.app.ListSignatures(ctx, filters, sort, pagination.Request{
		Page: req.Pag.Page,
		Size: req.Pag.Size,
	})
	if err != nil {
		logger.Log(ctx).Errorf("failed to list petition signers: %v", err)

		return nil, err
	}

	return responses.SignaturesList(signers, pag), nil
}
