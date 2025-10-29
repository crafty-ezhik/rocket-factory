package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

const (
	partsCollection = "parts"

	partFieldPartUUID            = "part_uuid"
	partFieldName                = "name"
	partFieldCategory            = "category"
	partFieldTags                = "tags"
	partFieldManufacturerCountry = "manufacturer.country"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) *repository {
	collection := db.Collection(partsCollection)

	// Добавляем индексы
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: partFieldPartUUID, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &repository{
		collection: collection,
	}
}
