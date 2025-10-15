package converter

import (
	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
)

// PartToServiceModel - преобразует модель репозитория в сервисную модель
func PartToServiceModel(part repoModel.Part) serviceModel.Part {
	return serviceModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category,
		Dimensions:    dimensionsToServiceModel(part.Dimensions),
		Manufacturer:  manufacturerToServiceModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// PartToRepoModel - Конвертирует сервисную модель в модель репозитория
func PartToRepoModel(part serviceModel.Part) repoModel.Part {
	return repoModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category,
		Dimensions:    dimensionsToRepoModel(part.Dimensions),
		Manufacturer:  manufacturerToRepoModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func dimensionsToServiceModel(dimensions *repoModel.Dimensions) *serviceModel.Dimensions {
	return &serviceModel.Dimensions{
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Length: dimensions.Length,
		Weight: dimensions.Weight,
	}
}

func dimensionsToRepoModel(dimensions *serviceModel.Dimensions) *repoModel.Dimensions {
	return &repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func manufacturerToServiceModel(manufacturer *repoModel.Manufacturer) *serviceModel.Manufacturer {
	return &serviceModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Website,
		Website: manufacturer.Website,
	}
}

func manufacturerToRepoModel(manufacturer *serviceModel.Manufacturer) *repoModel.Manufacturer {
	return &repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Website,
		Website: manufacturer.Website,
	}
}
