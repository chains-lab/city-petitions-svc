package app

import (
	"database/sql"

	"github.com/chains-lab/city-petitions-svc/internal/app/entities"
	"github.com/chains-lab/city-petitions-svc/internal/config"
)

type App struct {
	entities.Petition
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	return App{
		Petition: entities.NewPetition(pg),
	}, nil
}
