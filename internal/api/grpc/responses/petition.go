package responses

import (
	pagProto "github.com/chains-lab/city-petitions-proto/gen/go/common/pagination"
	svc "github.com/chains-lab/city-petitions-proto/gen/go/svc/petition"
	"github.com/chains-lab/city-petitions-svc/internal/app/models"
	"github.com/chains-lab/city-petitions-svc/internal/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Petition(model models.Petition) *svc.Petition {
	return &svc.Petition{
		Id:          model.ID.String(),
		CityId:      model.CityID.String(),
		Title:       model.Title,
		Description: model.Description,
		CreatorId:   model.CreatorID.String(),
		Status:      model.Status,
		Signatures:  uint32(model.Signatures),
		Goal:        uint32(model.Goal),
		Reply:       model.Reply,
		EndDate:     timestamppb.New(model.EndDate),
		CreatedAt:   timestamppb.New(model.CreatedAt),
		UpdatedAt:   timestamppb.New(model.UpdatedAt),
	}
}

func PetitionsList(models []models.Petition, pagResp pagination.Response) *svc.PetitionList {
	petitions := make([]*svc.Petition, 0, len(models))

	for _, model := range models {
		petitions = append(petitions, Petition(model))
	}

	return &svc.PetitionList{
		Petitions: petitions,
		Pagination: &pagProto.Response{
			Page:  pagResp.Page,
			Size:  pagResp.Size,
			Total: pagResp.Total,
		},
	}
}
