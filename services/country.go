package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"fmt"
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
