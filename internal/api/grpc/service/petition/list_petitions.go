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

func (s Service) ListPetitions(cxt context.Context, req *svc.ListPetitionsRequest) (*svc.PetitionList, error) {
	filters := entities.ListPetitionsFilter{}

	if req.Filters.CityId != "" {
		cityId, err := uuid.Parse(req.Filters.CityId)
		if err != nil {
			logger.Log(cxt).Errorf("failed to parse city id: %v", err)

			return nil, problems.InvalidArgumentError(cxt, "city_id is invalid", &errdetails.BadRequest_FieldViolation{
				Field:       "city_id",
				Description: "invalid UUID format for city ID",
			})
		}

		filters.CityID = &cityId
	}

	if req.Filters.CreatorId != "" {
		creatorId, err := uuid.Parse(req.Filters.CreatorId)
		if err != nil {
			logger.Log(cxt).Errorf("failed to parse creator id: %v", err)

			return nil, problems.InvalidArgumentError(cxt, "creator_id is invalid", &errdetails.BadRequest_FieldViolation{
				Field:       "creator_id",
				Description: "invalid UUID format for creator ID",
			})
		}

		filters.CreatorID = &creatorId
	}

	filters.TitleLike = &req.Filters.TitleLike
	filters.Rejected = &req.Filters.Rejected
	filters.Approved = &req.Filters.Approved
	filters.Available = &req.Filters.Available
	filters.Expired = &req.Filters.Expired

	sort := entities.ListPetitionsSort{}
	switch srt := req.Sort.(type) {
	case *svc.ListPetitionsRequest_Oldest:
		if srt.Oldest {
			sort.Oldest = true
		} else {
			sort.Oldest = true
		}
	case *svc.ListPetitionsRequest_LeastSignatures:
		if srt.LeastSignatures {
			sort.LessSign = true
		} else {
			sort.LessSign = true
		}
	case *svc.ListPetitionsRequest_MostSignatures:
		if srt.MostSignatures {
			sort.MoreSign = true
		} else {
			sort.MoreSign = true
		}
	default:
		sort.Newest = true
	}

	petitions, pag, err := s.app.ListPetitions(cxt, filters, sort, pagination.Request{
		Page: req.Pag.Page,
		Size: req.Pag.Size,
	})
	if err != nil {
		logger.Log(cxt).Errorf("failed to list petitions: %v", err)

		return nil, err
	}

	return responses.PetitionsList(petitions, pag), nil
}
