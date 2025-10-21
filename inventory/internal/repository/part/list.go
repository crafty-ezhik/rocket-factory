package part

import (
	"context"

	serviceModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/model"
)

type filter func(part repoModel.Part) bool

func (r *repository) List(_ context.Context, filters serviceModel.PartsFilter) ([]serviceModel.Part, error) {
	var filteredParts []repoModel.Part

	filtersFn := getFilterFuncs(filters)

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, part := range r.data {
		ok := true
		for _, filterFn := range filtersFn {
			if !filterFn(part) {
				ok = false
				break
			}
		}

		if ok {
			filteredParts = append(filteredParts, part)
		}
	}

	return converter.SlicePartToServiceModel(filteredParts), nil
}

// getFilterFuncs - возвращает список фильтрующих функций
func getFilterFuncs(filterValues serviceModel.PartsFilter) []filter {
	uuidSet := set(filterValues.UUIDs)
	nameSet := set(filterValues.Names)
	catSet := set(filterValues.Categories)
	countrySet := set(filterValues.ManufacturerCountry)
	tagSet := set(filterValues.Tags)

	return []filter{
		func(part repoModel.Part) bool {
			if len(uuidSet) > 0 {
				if _, ok := uuidSet[part.UUID.String()]; !ok {
					return false
				}
			}
			return true
		},
		func(part repoModel.Part) bool {
			if len(nameSet) > 0 {
				if _, ok := nameSet[part.Name]; !ok {
					return false
				}
			}
			return true
		},
		func(part repoModel.Part) bool {
			if len(catSet) > 0 {
				if _, ok := catSet[part.Category]; !ok {
					return false
				}
			}
			return true
		},
		func(part repoModel.Part) bool {
			if len(countrySet) > 0 {
				if _, ok := countrySet[part.Manufacturer.Country]; !ok {
					return false
				}
			}
			return true
		},
		func(part repoModel.Part) bool {
			if len(tagSet) > 0 {
				exist := false
				for _, partTag := range part.Tags {
					if _, ok := tagSet[partTag]; ok {
						exist = true
					}
				}
				if !exist {
					return false
				}
			}
			return true
		},
	}
}

// set - преобразует слайс в множество
func set(data []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range data {
		m[v] = struct{}{}
	}
	return m
}
