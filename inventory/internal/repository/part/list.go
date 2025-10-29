package part

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filters serviceModel.PartsFilter) ([]serviceModel.Part, error) {
	var filteredParts []repoModel.Part

	cursor, err := r.collection.Find(ctx, filtersToBson(filters))
	if err != nil {
		return nil, fmt.Errorf("error finding parts: %w", err)
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("error closing cursor: %v\n", cerr)
		}
	}()

	err = cursor.All(ctx, &filteredParts)
	if err != nil {
		return nil, fmt.Errorf("error getting parts: %w", err)
	}

	return converter.SlicePartToServiceModel(filteredParts), nil
}

func filtersToBson(filters serviceModel.PartsFilter) bson.M {
	if filters.IsEmpty() {
		return bson.M{}
	}

	var conditions bson.A
	if len(filters.UUIDs) > 0 {
		conditions = append(conditions, bson.M{
			partFieldPartUUID: bson.M{"$in": uuidToBsonA(filters.UUIDs)},
		})
	}
	if len(filters.Names) > 0 {
		if len(filters.Names) == 1 {
			conditions = append(conditions, bson.M{
				partFieldName: bson.M{"$regex": filters.Names[0], "$options": "i"},
			})
		} else {
			conditions = append(conditions, bson.M{
				partFieldName: bson.M{"$in": filters.Names},
			})
		}
	}
	if len(filters.Categories) > 0 {
		conditions = append(conditions, bson.M{
			partFieldCategory: bson.M{"$in": filters.Categories},
		})
	}
	if len(filters.ManufacturerCountry) > 0 {
		conditions = append(conditions, bson.M{
			partFieldManufacturerCountry: bson.M{"$in": filters.ManufacturerCountry},
		})
	}
	if len(filters.Tags) > 0 {
		conditions = append(conditions, bson.M{
			partFieldTags: bson.M{"$in": filters.Tags},
		})
	}

	return bson.M{"$and": conditions}
}

func uuidToBsonA(uuids []string) bson.A {
	arr := bson.A{}
	for _, v := range uuids {
		if uuidV, err := uuid.Parse(v); err == nil {
			arr = append(arr, primitive.Binary{Subtype: 0x00, Data: uuidV[:]})
		}
	}
	return arr
}
