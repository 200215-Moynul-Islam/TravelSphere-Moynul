package utils

import (
	"TravelSphere/constants"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestGetAllCountries(t *testing.T) {
	testCases := []struct {
		name string
		mockStatusCode int
		mockResponseBody string
		expectError bool
		overrideURL string
		clearConfig bool
	}{
		{
			name: "Success - Valid country array returned",
			mockStatusCode: http.StatusOK,
			mockResponseBody: `{"data": {"objects": [
				{"names": {"common": "Albania"}, "codes": {"alpha_3": "ALB"}, "population": 2800000},
				{"names": {"common": "Algeria"}, "codes": {"alpha_3": "DStandard"}, "population": 44000000}
			]}}`,
			expectError: false,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "API Error - Remote server returns 500 status",
			mockStatusCode: http.StatusInternalServerError,
			mockResponseBody: `{"message": "Internal Server Error"}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "JSON Parse Error - Server returns invalid JSON syntax",
			mockStatusCode: http.StatusOK,
			mockResponseBody: `{"data": {"objects": [{invalid-json}]}}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "Connection Error - Server URL protocol is completely invalid",
			mockStatusCode: http.StatusOK,
			mockResponseBody: ``,
			expectError: true,
			overrideURL: "http://invalid-url-domain-space-error.local",
			clearConfig: false,
		},
		{
			name: "Config Error Fallback - Triggers config missing branch and breaks on live fallback connection",
			mockStatusCode: http.StatusOK,
			mockResponseBody: ``,
			expectError: true,
			overrideURL: "http://invalid-fallback-live-url-route.local",
			clearConfig: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("response_fields") == "" {
					t.Errorf("Expected URL query parameters to contain response_fields selection target criteria")
				}
				if r.URL.Query().Get("limit") == "" {
					t.Errorf("Expected URL query parameters to contain limit constraint")
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.mockStatusCode)
				_, _ = w.Write([]byte(tc.mockResponseBody))
			}))
			defer mockServer.Close()

			targetURL := mockServer.URL
			if tc.overrideURL != "" {
				targetURL = tc.overrideURL
			}
			_ = beego.AppConfig.Set("restcountries_api_key", "mock-key")
			_ = beego.AppConfig.Set("restcountries_base_url", targetURL)

			if tc.clearConfig {
				_ = beego.AppConfig.Set("restcountries_base_url", "")
				if tc.overrideURL != "" {
					_ = beego.AppConfig.Set("restcountries_base_url", tc.overrideURL)
				}
			}

			resultData, err := GetAllCountries(constants.DefaultCountriesLimit)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an execution error, but received nil")
				}
				if resultData != nil {
					t.Errorf("Expected output payload slice to be nil on error states, got: %v", resultData)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no execution error, but got: %v", err)
			}

			var apiEnvelope struct {
				Data struct {
					Objects []CountryDTO `json:"objects"`
				} `json:"data"`
			}
			if unmarshalErr := json.Unmarshal([]byte(tc.mockResponseBody), &apiEnvelope); unmarshalErr != nil {
				t.Fatalf("Test setup failure while parsing mock JSON text string: %v", unmarshalErr)
			}
			expectedData := apiEnvelope.Data.Objects

			if !reflect.DeepEqual(resultData, expectedData) {
				t.Errorf("Returned data does not match the mock API data exactly.\nReturned: %+v\nExpected: %+v", resultData, expectedData)
			}
		})
	}
}

func TestGetCountriesByRegion(t *testing.T) {
	testCases := []struct {
		name string
		inputRegion string
		mockStatusCode int
		mockResponseBody string
		expectError bool
		overrideURL string
		clearConfig bool
	}{
		{
			name: "Success - Valid country array returned for region",
			inputRegion: "europe",
			mockStatusCode: http.StatusOK,
			mockResponseBody: `{"data": {"objects": [
				{"names": {"common": "Germany"}, "codes": {"alpha_3": "DEU"}, "region": "Europe"},
				{"names": {"common": "France"}, "codes": {"alpha_3": "FRA"}, "region": "Europe"}
			]}}`,
			expectError: false,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "API Error - Remote server returns 500 status",
			inputRegion: "asia",
			mockStatusCode: http.StatusInternalServerError,
			mockResponseBody: `{"message": "Internal Server Error"}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "JSON Parse Error - Server returns invalid JSON syntax",
			inputRegion: "africa",
			mockStatusCode: http.StatusOK,
			mockResponseBody: `{"data": {"objects": [{invalid-json}]}}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "Connection Error - Server URL protocol is completely invalid",
			inputRegion: "americas",
			mockStatusCode: http.StatusOK,
			mockResponseBody: ``,
			expectError: true,
			overrideURL: "http://invalid-url-domain-space-error.local",
			clearConfig: false,
		},
		{
			name: "Config Error Fallback - Triggers config missing branch and breaks on live fallback connection",
			inputRegion: "oceania",
			mockStatusCode: http.StatusOK,
			mockResponseBody: ``,
			expectError: true,
			overrideURL: "",
			clearConfig: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("region") != tc.inputRegion {
					t.Errorf("Expected URL region query parameter to be %s, but got: %s", tc.inputRegion, r.URL.Query().Get("region"))
				}
				if r.URL.Query().Get("response_fields") == "" {
					t.Errorf("Expected URL query parameters to contain response_fields selection target criteria")
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.mockStatusCode)
				_, _ = w.Write([]byte(tc.mockResponseBody))
			}))
			defer mockServer.Close()

			targetURL := mockServer.URL
			if tc.overrideURL != "" {
                targetURL = tc.overrideURL
			}
			_ = beego.AppConfig.Set("restcountries_api_key", "mock-key")
			_ = beego.AppConfig.Set("restcountries_base_url", targetURL)

			if tc.clearConfig {
				_ = beego.AppConfig.Set("restcountries_base_url", "")
				originalTransport := http.DefaultTransport
				defer func() { http.DefaultTransport = originalTransport }()
				http.DefaultTransport = &http.Transport{
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						return nil, fmt.Errorf("forced network failure for fallback testing")
					},
				}
			}

			resultData, err := GetCountriesByRegion(tc.inputRegion)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an execution error, but received nil")
				}
				if resultData != nil {
					t.Errorf("Expected output payload slice to be nil on error states, got: %v", resultData)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no execution error, but got: %v", err)
			}

			var apiEnvelope struct {
				Data struct {
					Objects []CountryDTO `json:"objects"`
				} `json:"data"`
			}
			if unmarshalErr := json.Unmarshal([]byte(tc.mockResponseBody), &apiEnvelope); unmarshalErr != nil {
				t.Fatalf("Test setup failure while parsing mock JSON text string: %v", unmarshalErr)
			}
			expectedData := apiEnvelope.Data.Objects

			if !reflect.DeepEqual(resultData, expectedData) {
				t.Errorf("Returned data does not match the mock API data exactly.\nReturned: %+v\nExpected: %+v", resultData, expectedData)
			}
		})
	}
}
