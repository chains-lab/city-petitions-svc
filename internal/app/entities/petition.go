package entities

import (
	"context"

	"github.com/chains-lab/city-petitions-svc/internal/dbx"
	"github.com/google/uuid"
)

type petitionsQ interface {
	New() dbx.PetitionsQ

	Insert(ctx context.Context, in dbx.Petition) error
	Get(ctx context.Context) (dbx.Petition, error)
	Select(ctx context.Context) ([]dbx.Petition, error)
	Update(ctx context.Context, in dbx.UpdatePetitionInput) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.PetitionsQ
	FilterCityID(cityID uuid.UUID) dbx.PetitionsQ
	FilterCreatorID(userID uuid.UUID) dbx.PetitionsQ
	FilterStatus(status string) dbx.PetitionsQ

	TitleLike(s string) dbx.PetitionsQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.PetitionsQ
}

type signaturesQ interface {
	New() dbx.PetitionSignaturesQ

	Insert(ctx context.Context, input dbx.PetitionSignature) error
	Get(ctx context.Context) (dbx.PetitionSignature, error)
	Select(ctx context.Context) ([]dbx.PetitionSignature, error)
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.PetitionSignaturesQ
	FilterPetitionID(petitionID uuid.UUID) dbx.PetitionSignaturesQ
	FilterUserID(userID uuid.UUID) dbx.PetitionSignaturesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.PetitionSignaturesQ
}
