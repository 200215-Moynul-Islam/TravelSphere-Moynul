package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

func GetCountriesByCodes(codes []string) ([]map[string]any, error) {
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

	var rawData []map[string]any
	if err := json.NewDecoder(apiResponse.Body).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return rawData, nil
}
