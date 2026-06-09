package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"fmt"
	"strings"
)

type CountryService struct{}

func (s *CountryService) GetCountriesByCodes(codes []string) ([]models.Country, error) {
	if len(codes) == 0 {
		return []models.Country{}, nil
	}

	dtos, err := utils.GetCountriesByCodes(codes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch countries: %w", err)
	}

	countries, err := utils.MapToCountrySlice(dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to translate country data: %w", err)
	}

	return countries, nil
}

func (s *CountryService) GetCountriesByPartialName(partialName string) ([]models.Country, error) {
    normalizedName := strings.TrimSpace(partialName)
    if normalizedName == "" {
        return []models.Country{}, nil
    }

    dtos, err := utils.GetCountriesByPartialName(normalizedName)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch countries: %w", err)
    }

    countries, err := utils.MapToCountrySlice(dtos)
    if err != nil {
        return nil, fmt.Errorf("failed to translate country data: %w", err)
    }

    return countries, nil
}
