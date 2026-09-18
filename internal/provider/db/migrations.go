package db

import (
	"context"
	"io/fs"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func runMigrations(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) error {
	fsys, err := fs.Sub(embeddedFS, "migrations")
	if err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)

	provider, err := goose.NewProvider(goose.DialectPostgres, db, fsys)
	if err != nil {
		return err
	}
	defer func() {
		if err := provider.Close(); err != nil {
			logger.Error("closing goose provider conn", "err", err)
		}
	}()

	results, err := provider.Up(ctx)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		logger.Info("no migrations applied")
	}

	for _, result := range results {
		logger.Info("migration applied", "info", result.String())
	}

	return nil
}
