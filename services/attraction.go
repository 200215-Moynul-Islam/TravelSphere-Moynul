package services

import (
	"TravelSphere/models"
	"TravelSphere/utils"
	"strings"
)

type AttractionService struct{}

func (s *AttractionService) GetPopularAttractions(cityName string) ([]models.Attraction, error) {
	if strings.TrimSpace(cityName) == "" {
		return []models.Attraction{}, nil
	}

	geo, err := utils.FetchCoordinatesByCity(cityName)
	if err != nil {
		return nil, err
	}

	features, err := utils.FetchAttractionsByCoordinates(geo.Lat, geo.Lon)
	if err != nil {
		return nil, err
	}

	var results []models.Attraction
	for _, f := range features.Features {
		if f.Properties.Name == "" {
			continue
		}

		rawTags := strings.Split(f.Properties.Kinds, ",")
		var cleanTags []string
		
		for _, tag := range rawTags {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				cleanTags = append(cleanTags, trimmed)
			}
		}

		results = append(results, models.Attraction{
			Name: f.Properties.Name,
			Tags: cleanTags,
		})
	}

	return results, nil
}
