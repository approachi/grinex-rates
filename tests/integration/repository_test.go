package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/approachi/grinex-rates/internal/domain"
	"github.com/approachi/grinex-rates/internal/repository/postgres"
)

func TestRepository_ShouldPersistRate_WhenSave(t *testing.T) {
	pool := newTestPool(t)
	truncateRates(t, pool)

	repo := postgres.NewRepository(pool)

	rate := domain.Rate{
		Ask:       8085,
		Bid:       7900,
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}

	err := repo.Save(context.Background(), rate)
	require.NoError(t, err)

	var ask, bid int64
	var createdAt time.Time
	err = pool.QueryRow(context.Background(),
		`SELECT ask, bid, created_at FROM rates LIMIT 1`,
	).Scan(&ask, &bid, &createdAt)

	require.NoError(t, err)
	assert.Equal(t, rate.Ask, ask)
	assert.Equal(t, rate.Bid, bid)
	assert.WithinDuration(t, rate.CreatedAt, createdAt, time.Second)
}

func TestRepository_ShouldPersistAllRates_WhenSaveMultiple(t *testing.T) {
	pool := newTestPool(t)
	truncateRates(t, pool)

	repo := postgres.NewRepository(pool)

	for i := range 3 {
		err := repo.Save(context.Background(), domain.Rate{
			Ask:       int64(8000 + i),
			Bid:       int64(7900 + i),
			CreatedAt: time.Now().UTC(),
		})
		require.NoError(t, err)
	}

	var count int
	err := pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM rates`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}
