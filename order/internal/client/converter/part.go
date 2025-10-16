package converter

import (
	"log"

	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	genInventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filters serviceModel.PartsFilter) *genInventoryV1.PartsFilter {
	categories := make([]genInventoryV1.Category, 0, len(filters.Categories))
	for _, category := range filters.Categories {
		categories = append(categories, genInventoryV1.Category(genInventoryV1.Category_value[category]))
	}

	return &genInventoryV1.PartsFilter{
		Uuids:                 filters.UUIDs,
		Names:                 filters.Names,
		Categories:            categories,
		ManufacturerCountries: filters.ManufacturerCountry,
		Tags:                  filters.Tags,
	}
}

func PartListToServiceModel(parts []*genInventoryV1.Part) []serviceModel.Part {
	result := make([]serviceModel.Part, len(parts))
	for i, part := range parts {
		result[i] = PartToServiceModel(part)
	}
	return result
}

func PartToServiceModel(part *genInventoryV1.Part) serviceModel.Part {
	return serviceModel.Part{
		UUID:          uuid.MustParse(part.Uuid),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category.String(),
		Dimensions:    dimensionsToServiceModel(part.Dimensions),
		Manufacturer:  manufacturerToServiceModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataToServiceModel(part.Metadata),
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     part.UpdatedAt.AsTime(),
	}
}

func dimensionsToServiceModel(dim *genInventoryV1.Dimensions) *serviceModel.Dimensions {
	return &serviceModel.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func manufacturerToServiceModel(man *genInventoryV1.Manufacturer) *serviceModel.Manufacturer {
	return &serviceModel.Manufacturer{
		Name:    man.Name,
		Country: man.Country,
		Website: man.Website,
	}
}

func metadataToServiceModel(md map[string]*genInventoryV1.Value) map[string]any {
	result := make(map[string]any, len(md))
	for k, v := range md {
		if v == nil || v.GetValueType() == nil {
			continue
		}
		switch v.GetValueType().(type) {
		case *genInventoryV1.Value_StringValue:
			result[k] = v.GetStringValue()
		case *genInventoryV1.Value_BoolValue:
			result[k] = v.GetBoolValue()
		case *genInventoryV1.Value_DoubleValue:
			result[k] = v.GetDoubleValue()
		case *genInventoryV1.Value_Int64Value:
			result[k] = v.GetInt64Value()
		default:
			log.Printf("unknown value type: %T", v.GetValueType())
		}
	}
	return result
}
