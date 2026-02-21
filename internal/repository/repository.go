package repository

import (
	"context"
	"device-parser-logs/internal/models"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	DeviceLogs    *mongo.Collection
	ProcessedFile *mongo.Collection
}

func New(ctx context.Context, dbName string, client *mongo.Client) *MongoDB {
	return &MongoDB{
		DeviceLogs:    client.Database(dbName).Collection("DeviceLogs"),
		ProcessedFile: client.Database(dbName).Collection("ProcessedFile"),
	}

}

func (p *MongoDB) SaveDeviceLogs(ctx context.Context, data []*models.DeviceLogs) error {
	res, err := p.DeviceLogs.InsertMany(ctx, data)
	if !res.Acknowledged {
		return fmt.Errorf("writing was unacknowledged: %w", err)
	}
	return err
}

func (p *MongoDB) GetInfoByGuid(ctx context.Context, guid string, page, limit int) (*models.PaginationResult, error) {
	filter := bson.M{"guid": guid}

	total, err := p.DeviceLogs.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := p.DeviceLogs.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find units: %w", err)
	}
	defer cursor.Close(ctx)

	var deviceLogs []*models.DeviceLogs
	if err = cursor.All(ctx, &deviceLogs); err != nil {
		return nil, fmt.Errorf("failed to decode deviceLogs: %w", err)
	}

	result := &models.PaginationResult{
		Data:  deviceLogs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return result, nil
}

func (p *MongoDB) GetProssecedFile(ctx context.Context, fileName string) (*models.ProcessedFile, error) {
	filter := bson.M{"name": fileName}

	var file models.ProcessedFile

	err := p.ProcessedFile.FindOne(ctx, filter).Decode(&file)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get file by filename: %w", err)
	}

	return &file, nil
}

func (p *MongoDB) SaveProssecedFile(ctx context.Context, prossecedFile *models.ProcessedFile) error {
	prossecedFile.ID = primitive.NewObjectID()
	_, err := p.ProcessedFile.InsertOne(ctx, prossecedFile)
	return err
}
