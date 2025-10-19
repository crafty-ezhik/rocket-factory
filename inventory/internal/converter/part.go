package converter

import (
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

// SlicePartToProto - конвертация []serviceModel.Part в []*inventoryV1.Part
func SlicePartToProto(parts []serviceModel.Part) []*inventoryV1.Part {
	protoParts := make([]*inventoryV1.Part, len(parts))
	for i, part := range parts {
		protoParts[i] = PartToProto(part)
	}
	return protoParts
}

// PartToProto - Конвертация serviceModel в inventoryV1.Part
func PartToProto(part serviceModel.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      inventoryV1.Category(inventoryV1.Category_value[part.Category]),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      convertMapToValues(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
	}
}

// PartsFilterToServiceModel - Конвертация inventoryV1.PartFilters в serviceModel
func PartsFilterToServiceModel(filters *inventoryV1.PartsFilter) serviceModel.PartsFilter {
	if filters == nil {
		return serviceModel.PartsFilter{}
	}
	return serviceModel.PartsFilter{
		UUIDs:               filters.GetUuids(),
		Names:               filters.GetNames(),
		Categories:          cat2string(filters.GetCategories()),
		ManufacturerCountry: filters.GetManufacturerCountries(),
		Tags:                filters.GetTags(),
	}
}

// dimensionsToProto - Конвертация serviceModel.Dimensions в inventoryV1.Dimensions
func dimensionsToProto(dim *serviceModel.Dimensions) *inventoryV1.Dimensions {
	if dim == nil {
		return &inventoryV1.Dimensions{}
	}
	return &inventoryV1.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

// Конвертация serviceModel.Manufacturer в inventoryV1.Manufacturer
func manufacturerToProto(man *serviceModel.Manufacturer) *inventoryV1.Manufacturer {
	if man == nil {
		return &inventoryV1.Manufacturer{}
	}
	return &inventoryV1.Manufacturer{
		Name:    man.Name,
		Country: man.Country,
		Website: man.Website,
	}
}

// cat2string - конвертирует тип []inventoryV1.Category в []string для фильтрации
func cat2string(cats []inventoryV1.Category) []string {
	result := make([]string, 0, len(cats))
	for _, v := range cats {
		result = append(result, v.String())
	}
	return result
}

// convertMapToValues конвертирует map[string]any в map[string]*inventoryV1.Value.
func convertMapToValues(input map[string]any) map[string]*inventoryV1.Value {
	if input == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value, len(input))

	for key, val := range input {
		value := &inventoryV1.Value{}

		switch v := val.(type) {
		case string:
			value.ValueType = &inventoryV1.Value_StringValue{StringValue: v}
		case bool:
			value.ValueType = &inventoryV1.Value_BoolValue{BoolValue: v}
		case int64:
			value.ValueType = &inventoryV1.Value_Int64Value{Int64Value: v}
		case int:
			value.ValueType = &inventoryV1.Value_Int64Value{Int64Value: int64(v)}
		case int32:
			value.ValueType = &inventoryV1.Value_Int64Value{Int64Value: int64(v)}
		case float64:
			value.ValueType = &inventoryV1.Value_DoubleValue{DoubleValue: v}
		case float32:
			value.ValueType = &inventoryV1.Value_DoubleValue{DoubleValue: float64(v)}
		default:
			value.ValueType = &inventoryV1.Value_StringValue{StringValue: fmt.Sprintf("%v", v)}
		}
		result[key] = value
	}

	return result
}
