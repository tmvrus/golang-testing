package storage

import "context"

type Payout interface {
	Register(ctx context.Context, userID string, reqID int64, payout float64) error
	Count(ctx context.Context, userID string) (int64, error)
}
