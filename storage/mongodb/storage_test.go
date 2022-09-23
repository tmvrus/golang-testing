package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Storage(t *testing.T) {
	t.Skip("skip")

	time.Sleep(time.Minute * 10)

	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		const userID = "632daa2d1a2654ad2350b3e5"

		err := payoutStorage.Register(ctx, userID, 100, 1.1)
		require.NoError(t, err)
		err = payoutStorage.Register(ctx, userID, 101, 2.2)
		require.NoError(t, err)
		err = payoutStorage.Register(ctx, userID, 102, 3.3)
		require.NoError(t, err)

		count, err := payoutStorage.Count(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, count, int64(3))
	})

	t.Run("empty result for new user", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		const userID = "632daa2d1a2654ad2350b3ee"

		count, err := payoutStorage.Count(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, count, int64(0))
	})

	t.Run("idempotency", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		const (
			userID    = "632daa2d1a2654ad2350b3dd"
			requestID = 10
		)

		err := payoutStorage.Register(ctx, userID, requestID, 11)
		require.NoError(t, err)
		err = payoutStorage.Register(ctx, userID, requestID, 11)
		require.NoError(t, err)
		err = payoutStorage.Register(ctx, userID, requestID, 11)
		require.NoError(t, err)

		count, err := payoutStorage.Count(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, count, int64(1))
	})
}
