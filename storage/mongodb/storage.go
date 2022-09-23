package mongodb

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	dbName = "users"

	payoutCollectionName = "payouts"

	duplicateKeyErrorCode = "11000"
)

type Storage struct {
	db *mongo.Database
}

func NewStorage(client *mongo.Client) (*Storage, error) {
	ctx := context.Background()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	return &Storage{db: client.Database(dbName)}, nil
}

func (s *Storage) Count(ctx context.Context, userID string) (int64, error) {
	res, err := s.
		db.
		Collection(payoutCollectionName).
		CountDocuments(ctx, bson.M{"user_id": userID})

	if err != nil {
		return 0, fmt.Errorf("failed to count payouts: %w", err)
	}

	return res, nil
}

func (s *Storage) Register(ctx context.Context, userID string, reqID int64, payout float64) error {
	_, err := s.
		db.
		Collection(payoutCollectionName).
		InsertOne(
			ctx,
			bson.M{
				"request_id": reqID,
				"user_id":    userID,
				"payout":     payout,
			},
		)

	if err != nil {
		if strings.Contains(err.Error(), duplicateKeyErrorCode) {
			return nil
		}
		return err
	}

	return nil
}
