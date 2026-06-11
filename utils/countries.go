package utils

import (
	"TravelSphere/constants"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

func GetCountriesByCodes(codes []string) ([]CountryDTO, error) {
	endpointPath := fmt.Sprintf("?codes=%s", strings.Join(codes, ","))
	return callRestCountriesAPI(http.MethodGet, endpointPath)
}

func GetAllCountries() ([]CountryDTO, error) {
	endpointPath := fmt.Sprintf("?response_fields=%s", constants.RestCountriesFields)
	return callRestCountriesAPI(http.MethodGet, endpointPath)
}

func GetCountriesByRegion(region string) ([]CountryDTO, error) {
	endpointPath := fmt.Sprintf("?region=%s&response_fields=%s", region, constants.RestCountriesFields)
	return callRestCountriesAPI(http.MethodGet, endpointPath)
}

func callRestCountriesAPI(httpMethod string, endpointPath string) ([]CountryDTO, error) {
	apiKey, err := beego.AppConfig.String("restcountries_api_key")
	if err != nil || apiKey == "" {
		return nil, fmt.Errorf("restcountries configuration key missing from app config")
	}
	baseURL, err := beego.AppConfig.String("restcountries_base_url")
	if err != nil || baseURL == "" {
		baseURL = "https://api.restcountries.com/countries/v5"
	}
	url := fmt.Sprintf("%s%s", baseURL, endpointPath)
	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	apiResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed connection: %w", err)
	}
	defer apiResponse.Body.Close()

	if apiResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error status: %d", apiResponse.StatusCode)
	}
	var apiEnvelope struct {
		Data struct {
			Objects []CountryDTO `json:"objects"`
		} `json:"data"`
	}
	if err := json.NewDecoder(apiResponse.Body).Decode(&apiEnvelope); err != nil {
		return nil, fmt.Errorf("failed to decode wrapped json stream: %w", err)
	}
	return apiEnvelope.Data.Objects, nil
}
