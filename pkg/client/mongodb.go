package client

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password string) (*mongo.Client, error) {
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port)

	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoURI).SetAuth(options.Credential{
			Username: username,
			Password: password,
		}))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping to mongoDB: %w", err)
	}

	return client, nil
}

func CloseMongoDB(client *mongo.Client) error {
	if err := client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("failed to disconnect to mongoDB: %w", err)
	}
	return nil
}
