package petition

import (
	"context"

	svc "github.com/chains-lab/city-petitions-proto/gen/go/svc/petition"
	"github.com/chains-lab/city-petitions-svc/internal/app/entities"
	"github.com/chains-lab/city-petitions-svc/internal/app/models"
	"github.com/chains-lab/city-petitions-svc/internal/pagination"
	"github.com/google/uuid"

	"github.com/chains-lab/city-petitions-svc/internal/app"
	"github.com/chains-lab/city-petitions-svc/internal/config"
)

type application interface {
	CreatePetition(ctx context.Context, cityID, creatorID uuid.UUID, input entities.CreatePetitionInput) (models.Petition, error)
	GetPetition(ctx context.Context, petitionID uuid.UUID) (models.Petition, error)
	ApprovePetition(ctx context.Context, petitionID uuid.UUID, reply string) (models.Petition, error)
	RejectPetition(ctx context.Context, petitionID uuid.UUID, reply string) (models.Petition, error)

	SignPetition(ctx context.Context, initiatorID, petitionID uuid.UUID) (models.PetitionSignature, error)
	GetSignatureByID(ctx context.Context, userID, petitionID uuid.UUID) (models.PetitionSignature, error)

	GetSignatureByUserIDAndSigID(ctx context.Context, sigID uuid.UUID) (models.PetitionSignature, error)

	ListPetitions(
		ctx context.Context,
		filter entities.ListPetitionsFilter,
		sort entities.ListPetitionsSort,
		pag pagination.Request,
	) ([]models.Petition, pagination.Response, error)

	ListSignatures(
		ctx context.Context,
		filter entities.ListPetitionsSignFilter,
		sort entities.ListPetitionsSort,
		pag pagination.Request,
	) ([]models.PetitionSignature, pagination.Response, error)
}

type Service struct {
	app *app.App //TODO change to app.Application when refactored
	cfg config.Config

	svc.UnimplementedPetitionServiceServer
}

func NewService(cfg config.Config, app *app.App) Service {
	return Service{
		app: app,
		cfg: cfg,
	}
}
