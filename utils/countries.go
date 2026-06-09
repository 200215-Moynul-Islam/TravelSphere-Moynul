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
	baseURL, err := beego.AppConfig.String("restcountries_base_url")
	if err != nil || baseURL == "" {
		baseURL = "https://restcountries.com/v3.1" // Fallback url
	}
	url := fmt.Sprintf("%s/alpha?codes=%s", baseURL, strings.Join(codes, ","))
	apiResponse, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed connection: %w", err)
	}
	defer apiResponse.Body.Close()

	if apiResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error status: %d", apiResponse.StatusCode)
	}

	var countries []CountryDTO
	if err := json.NewDecoder(apiResponse.Body).Decode(&countries); err != nil {
		return nil, fmt.Errorf("failed to decode json stream: %w", err)
	}

	return countries, nil
}

func GetCountriesByPartialName(partialName string) ([]CountryDTO, error) {
    baseURL, err := beego.AppConfig.String("restcountries_base_url")
    if err != nil || baseURL == "" {
        baseURL = "https://restcountries.com/v3.1"
    }
    
    url := fmt.Sprintf("%s/name/%s?fullText=false", baseURL, partialName)
    apiResponse, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed connection: %w", err)
    }
    defer apiResponse.Body.Close()

    if apiResponse.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("api error status: %d", apiResponse.StatusCode)
    }

    var countries []CountryDTO
    if err := json.NewDecoder(apiResponse.Body).Decode(&countries); err != nil {
        return nil, fmt.Errorf("failed to decode json stream: %w", err)
    }

    return countries, nil
}

func GetAllCountries() ([]CountryDTO, error) {
	baseURL, err := beego.AppConfig.String("restcountries_base_url")
	if err != nil || baseURL == "" {
		baseURL = "https://restcountries.com/v3.1"
	}

	url := fmt.Sprintf("%s/all?fields=%s", baseURL, constants.RestCountriesFields)
	apiResponse, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed connection: %w", err)
	}
	defer apiResponse.Body.Close()

	if apiResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error status: %d", apiResponse.StatusCode)
	}

	var countries []CountryDTO
	if err := json.NewDecoder(apiResponse.Body).Decode(&countries); err != nil {
		return nil, fmt.Errorf("failed to decode json stream: %w", err)
	}

	return countries, nil
}

func GetCountriesByRegion(region string) ([]CountryDTO, error) {
    baseURL, err := beego.AppConfig.String("restcountries_base_url")
    if err != nil || baseURL == "" {
        baseURL = "https://restcountries.com/v3.1"
    }

    url := fmt.Sprintf("%s/region/%s?fields=%s", baseURL, region, constants.RestCountriesFields)
    apiResponse, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed connection: %w", err)
    }
    defer apiResponse.Body.Close()

    if apiResponse.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("api error status: %d", apiResponse.StatusCode)
    }

    var countries []CountryDTO
    if err := json.NewDecoder(apiResponse.Body).Decode(&countries); err != nil {
        return nil, fmt.Errorf("failed to decode json stream: %w", err)
    }

    return countries, nil
}
