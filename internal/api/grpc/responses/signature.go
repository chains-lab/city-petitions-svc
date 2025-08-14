package responses

import (
	pagProto "github.com/chains-lab/city-petitions-proto/gen/go/common/pagination"
	svc "github.com/chains-lab/city-petitions-proto/gen/go/svc/petition"
	"github.com/chains-lab/city-petitions-svc/internal/app/models"
	"github.com/chains-lab/city-petitions-svc/internal/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Signature(model models.PetitionSignature) *svc.Signature {
	return &svc.Signature{
		Id:         model.ID.String(),
		PetitionId: model.PetitionID.String(),
		UserId:     model.UserID.String(),
		CreatedAt:  timestamppb.New(model.CreatedAt),
	}
}

func SignaturesList(models []models.PetitionSignature, pagResp pagination.Response) *svc.SignatureList {
	signatures := make([]*svc.Signature, 0, len(models))

	for _, model := range models {
		signatures = append(signatures, Signature(model))
	}

	return &svc.SignatureList{
		Signatures: signatures,
		Pagination: &pagProto.Response{
			Page:  pagResp.Page,
			Size:  pagResp.Size,
			Total: pagResp.Total,
		},
	}
}
