package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	beego "github.com/beego/beego/v2/server/web"
)

type GeonameDTO struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type AttractionFeaturesDTO struct {
	Features []struct {
		Properties struct {
			Name  string `json:"name"`
			Kinds string `json:"kinds"`
		} `json:"properties"`
	} `json:"features"`
}

func FetchCoordinatesByCity(cityName string) (*GeonameDTO, error) {
	apiKey, err := beego.AppConfig.String("opentripmap_api_key")
	if err != nil || apiKey == "" {
		return nil, fmt.Errorf("opentripmap configuration key missing from app config")
	}

	baseURL, err := beego.AppConfig.String("opentripmap_geoname_url")
	if err != nil || baseURL == "" {
		baseURL = "https://api.opentripmap.com/0.1/en/places/geoname"
	}

	apiURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	queryParams := apiURL.Query()
	queryParams.Set("name", cityName)
	queryParams.Set("apikey", apiKey)
	apiURL.RawQuery = queryParams.Encode()

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return nil, fmt.Errorf("network failure during geoname lookup: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opentripmap geoname api returned bad status code: %d", resp.StatusCode)
	}

	var data GeonameDTO
	if decodeErr := json.NewDecoder(resp.Body).Decode(&data); decodeErr != nil {
		return nil, fmt.Errorf("failed parsing geoname json body payload: %w", decodeErr)
	}

	return &data, nil
}

func FetchAttractionsByCoordinates(lat, lon float64) (*AttractionFeaturesDTO, error) {
	apiKey, err := beego.AppConfig.String("opentripmap_api_key")
	if err != nil || apiKey == "" {
		return nil, fmt.Errorf("opentripmap configuration key missing from app config")
	}

	baseURL, err := beego.AppConfig.String("opentripmap_radius_url")
	if err != nil || baseURL == "" {
		baseURL = "https://api.opentripmap.com/0.1/en/places/radius"
	}

	apiURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	queryParams := apiURL.Query()
	queryParams.Set("radius", "10000")
	queryParams.Set("lon", fmt.Sprintf("%f", lon))
	queryParams.Set("lat", fmt.Sprintf("%f", lat))
	queryParams.Set("kinds", "historic,architecture,cultural")
	queryParams.Set("limit", "5")
	queryParams.Set("apikey", apiKey)
	apiURL.RawQuery = queryParams.Encode()

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return nil, fmt.Errorf("network failure during radius lookup: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opentripmap radius api returned bad status code: %d", resp.StatusCode)
	}

	var data AttractionFeaturesDTO
	if decodeErr := json.NewDecoder(resp.Body).Decode(&data); decodeErr != nil {
		return nil, fmt.Errorf("failed parsing radius features json body payload: %w", decodeErr)
	}

	return &data, nil
}
