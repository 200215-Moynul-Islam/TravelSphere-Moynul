package services

import (
	"TravelSphere/constants"
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

    countryDTOs, err := utils.GetAllCountries(constants.DefaultCountriesLimit)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch countries: %w", err)
    }

    // Prepare a map for search codes for quick lookups
    targetCodes := make(map[string]bool, len(codes))
    for _, code := range codes {
        targetCodes[strings.TrimSpace(code)] = true
    }

    // Filter the coutries based on provied codes
    var matchedCountryDTOs []utils.CountryDTO
    for _, dto := range countryDTOs {
        if targetCodes[strings.ToUpper(dto.Codes.Alpha3)] {
            matchedCountryDTOs = append(matchedCountryDTOs, dto)
        }
    }

    filteredCountries, err := utils.MapToCountrySlice(matchedCountryDTOs)
    if err != nil {
        return nil, fmt.Errorf("failed to translate country data: %w", err)
    }

    return filteredCountries, nil
}

func (s *CountryService) GetAllCountries(limit int) ([]models.Country, error) {
    dtos, err := utils.GetAllCountries(limit)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch all countries: %w", err)
    }

    countries, err := utils.MapToCountrySlice(dtos)
    if err != nil {
        return nil, fmt.Errorf("failed to translate all country data: %w", err)
    }

    return countries, nil
}

// Returns a list of countries filtered by search term and/or region.
// If search is empty, returns all countries in the specified region (or all regions if "all" is specified).
// If search is provided, filters the region results further by matching the search term against country names and capitals.
func (s *CountryService) GetFilteredCountries(search, region string) ([]models.Country, error) {
    var dtos []utils.CountryDTO
    var err error

    region = strings.TrimSpace(strings.ToLower(region))
    search = strings.TrimSpace(strings.ToLower(search))

    if region != "" && region != "all" {
        dtos, err = utils.GetCountriesByRegion(region)
    } else {
        dtos, err = utils.GetAllCountries(constants.DefaultCountriesLimit)
    }

    if err != nil {
        return nil, err
    }
	
    countries, err := utils.MapToCountrySlice(dtos)
    if err != nil {
        return nil, fmt.Errorf("failed to translate filtered data: %w", err)
    }

    // Return the region-filtered slice immediately when search string is empty
    if search == "" {
        return countries, nil
    }

    // Filter countries by search term in name or capital (case-insensitive)
    var filtered []models.Country
    for _, country := range countries {
        nameMatch := strings.Contains(strings.ToLower(country.Name), search)
        capitalMatch := strings.Contains(strings.ToLower(country.Capital), search)

        if nameMatch || capitalMatch {
            filtered = append(filtered, country)
        }
    }

    return filtered, nil
}
