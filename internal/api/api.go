package api

import (
	"context"

	"github.com/chains-lab/city-petitions-svc/internal/api/grpc"
	"github.com/chains-lab/city-petitions-svc/internal/app"
	"github.com/chains-lab/city-petitions-svc/internal/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Start(ctx context.Context, cfg config.Config, log *logrus.Logger, app *app.App) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return grpc.Run(ctx, cfg, log, app) })

	return eg.Wait()
}
