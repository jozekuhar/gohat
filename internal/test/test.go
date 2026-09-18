package test

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// NewPostgresContainer runs "postgres:16" container and returns URI. It also returns
// terminate function which you need to defer immediately after creation.
func NewPostgresContainer(ctx context.Context) (string, func(), error) {
	dbName := "testdb"
	dbUsername := "testdb"
	dbPassword := "test1234"

	container, err := postgres.Run(
		ctx,
		"postgres:17",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUsername),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		return "", nil, err
	}

	terminate := func() {
		err := container.Terminate(ctx)
		if err != nil {
			log.Fatalf("error terminating container")
		}
	}

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		return "", nil, err
	}

	return uri, terminate, nil
}
