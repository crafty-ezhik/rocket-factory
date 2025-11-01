package part

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

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

	return &repository{
		collection: collection,
	}
}
