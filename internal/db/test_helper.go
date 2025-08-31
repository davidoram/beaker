package db

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestPostgresPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	postgresUrl := "postgres://postgres:password@localhost:5432/beaker_test?sslmode=disable"
	config, err := pgxpool.ParseConfig(postgresUrl)
	require.NoError(t, err)

	// Connect and create the pool
	pool, err := pgxpool.NewWithConfig(t.Context(), config)
	require.NoError(t, err)

	err = pool.Ping(t.Context())
	require.NoError(t, err, "Check the test db exists. Run `DB_ENV=test make recreate-db` to create it.")

	return pool
}
