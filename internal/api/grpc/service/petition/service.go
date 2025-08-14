package petition

import (
	svc "github.com/chains-lab/city-petitions-proto/gen/go/svc/petition"

	"github.com/chains-lab/city-petitions-svc/internal/app"
	"github.com/chains-lab/city-petitions-svc/internal/config"
)

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
