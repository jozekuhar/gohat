package auth

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/test"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	uri, terminate, err := test.NewPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("new postgres container: %v", err)
		return
	}
	defer terminate()

	testDB, err = db.ConnectAndMigrate(ctx, uri, slog.Default(), false)
	if err != nil {
		log.Fatalf("connect and migrate test db: %s", err)
		return
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}
