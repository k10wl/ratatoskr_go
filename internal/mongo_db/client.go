package mongo_db

import (
	"context"
	"ratatoskr/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client              *mongo.Client
	db                  *mongo.Database
	tagsCollection      *mongo.Collection
	analyticsCollection *mongo.Collection
}

func NewMongoDB(ctx context.Context, URI string, database string) (*MongoDB, error) {
	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(URI),
	)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	db := client.Database(database)
	return &MongoDB{
		client:              client,
		db:                  db,
		tagsCollection:      db.Collection("tags_menus"),
		analyticsCollection: db.Collection("tags_usage_statistics"),
	}, nil
}

func (m MongoDB) GetAllGroupsWithTags(ctx context.Context) (*[]models.Group, error) {
	c, err := m.tagsCollection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	var res []models.Group
	err = c.All(ctx, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
