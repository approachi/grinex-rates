//go:build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func dbURL() string {
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	return "host=localhost port=5433 user=postgres password=postgres dbname=grinex_test sslmode=disable"
}

func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), dbURL())
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })
	return pool
}

func truncateRates(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), "TRUNCATE TABLE rates")
	require.NoError(t, err)
}
