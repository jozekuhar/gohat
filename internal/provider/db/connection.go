package db

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations
var embeddedFS embed.FS

func ConnectAndMigrate(
	ctx context.Context,
	connString string,
	logger *slog.Logger,
	debug bool,
) (*pgxpool.Pool, error) {
	pool, err := connectPool(ctx, connString, logger)
	if err != nil {
		return nil, err
	}

	// We don't do migrations on debug so we avoid constant
	// and unfinished migrations in development.
	if debug {
		return pool, nil
	}

	err = runMigrations(ctx, pool, logger)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func connectPool(
	ctx context.Context,
	connString string,
	logger *slog.Logger,
) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db pool: %w", err)
	}

	return pool, nil
}
