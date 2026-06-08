package utils

import (
	"TravelSphere/models"
	"fmt"
	"strings"
)

func MapToCountryModel(dto CountryDTO) models.Country {
	capital := ""
	if len(dto.Capital) > 0 {
		capital = dto.Capital[0]
	}

	currencyStr := ""
	for code, curr := range dto.Currencies {
		currencyStr = fmt.Sprintf("%s (%s)", code, curr.Name)
		break
	}

	langList := make([]string, 0, len(dto.Languages))
	for _, lang := range dto.Languages {
		langList = append(langList, lang)
	}
	languagesStr := strings.Join(langList, ", ")

	regionStr := dto.Region
	if dto.Subregion != "" {
		regionStr = fmt.Sprintf("%s - %s", dto.Region, dto.Subregion)
	}

	return models.Country{
		Code: dto.Cca3,
		Name: dto.Name.Common,
		OfficialName: dto.Name.Official,
		Flag: dto.Flags.Png,
		Capital: capital,
		Population: formatPopulation(dto.Population),
		Region: regionStr,
		Currency: currencyStr,
		Languages: languagesStr,
	}
}

func MapToCountrySlice(dtos []CountryDTO) ([]models.Country, error) {
	countries := make([]models.Country, len(dtos))
	for i, dto := range dtos {
		countries[i] = MapToCountryModel(dto)
	}
	return countries, nil
}

func formatPopulation(pop int64) string {
	if pop >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(pop)/1000000)
	}
	if pop >= 1000 {
		return fmt.Sprintf("%.1fK", float64(pop)/1000)
	}
	return fmt.Sprintf("%d", pop)
}
