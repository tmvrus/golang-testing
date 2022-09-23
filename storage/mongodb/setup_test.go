package mongodb_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tmvrus/golang-testing/storage"
	"github.com/tmvrus/golang-testing/storage/mongodb"
)

var payoutStorage storage.Payout

func TestMain(m *testing.M) {
	shutdown, err := setup()
	if err != nil {
		log.Fatalf("failed to setup mongo tests: %v", err)
	}

	code := m.Run()

	shutdown()
	os.Exit(code)
}

func setup() (func(), error) {
	mongoContainer, err := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "mongo:4.2.10",
				ExposedPorts: []string{"27017"},
				Cmd:          []string{"--replSet", "rs0"},
				WaitingFor:   wait.ForHTTP("/").WithPort("27017/tcp"),
			},
			Started: true,
		},
	)
	if err != nil {
		return nil, err
	}
	cleanup := func() {
		_ = mongoContainer.Terminate(context.Background())
	}

	code, err := mongoContainer.Exec(context.Background(),
		[]string{"mongo", "--eval", `rs.initiate({_id: "rs0", members: [{_id: 0, host: "localhost:27017"}]})`},
	)
	if err != nil || code != 0 {
		cleanup()
		return nil, fmt.Errorf("failed to init mongo replica set (exit code %d): %w", code, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = wait.ForLog("transition to primary complete").WaitUntilReady(ctx, mongoContainer)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to wait for mongo rs setup")
	}

	uri, err := mongoContainer.PortEndpoint(context.Background(), "27017", "mongodb")
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to get mongo URI: %w", err)
	}

	opts := options.Client().
		ApplyURI(uri).
		SetReplicaSet("rs0").
		SetDirect(true)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	const (
		dbName               = "users"
		payoutCollectionName = "payouts"
	)

	coll := client.Database(dbName).Collection(payoutCollectionName)

	_, err = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{"request_id", 1},
		},
		Options: options.Index().
			SetName("uniq_request_id").
			SetUnique(true),
	})

	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	payoutStorage, err = mongodb.NewStorage(client)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	return cleanup, nil
}
