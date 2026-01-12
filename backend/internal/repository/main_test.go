package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx,
		"postgres://postgres:12345@localhost:5432/cvwo_test",
	)
	if err != nil {
		panic(err)
	}

	testDB = pool

	code := m.Run()

	pool.Close()
	os.Exit(code)
}
