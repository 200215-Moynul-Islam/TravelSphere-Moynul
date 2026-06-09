package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestFetchCoordinatesByCity(t *testing.T) {
	testCases := []struct {
		name string
		inputCity string
		mockStatusCode int
		mockResponseBody string
		expectError bool
		overrideURL string
		clearConfig bool
	}{
		{
			name: "Success - Valid coordinates returned",
			inputCity: "Paris",
			mockStatusCode: http.StatusOK,
			mockResponseBody: `{"lat": 48.8566, "lon": 2.3522}`,
			expectError: false,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "API Error - Remote server returns 500 status",
			inputCity: "Paris",
			mockStatusCode: http.StatusInternalServerError,
			mockResponseBody: `{"message": "Internal Server Error"}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "JSON Parse Error - Server returns invalid JSON syntax",
			inputCity: "Paris",
			mockStatusCode: http.StatusOK,
			mockResponseBody: `{"lat": 48.8566, invalid-json}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "Connection Error - Server URL protocol is invalid",
			inputCity: "Paris",
			mockStatusCode: http.StatusOK,
			mockResponseBody: ``,
			expectError: true,
			overrideURL: "http://invalid-url-domain-space-error.local",
			clearConfig: false,
		},
		{
			name: "Config Error - API key missing from configuration",
			inputCity: "Paris",
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
				if r.URL.Query().Get("name") != tc.inputCity {
					t.Errorf("Expected URL query name to be %s, but got: %s", tc.inputCity, r.URL.Query().Get("name"))
				}
				if r.URL.Query().Get("apikey") != "mock_key" {
					t.Errorf("Expected URL query apikey to be mock_key, but got: %s", r.URL.Query().Get("apikey"))
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

			_ = beego.AppConfig.Set("opentripmap_api_key", "mock_key")
			_ = beego.AppConfig.Set("opentripmap_geoname_url", targetURL)

			if tc.clearConfig {
				_ = beego.AppConfig.Set("opentripmap_api_key", "")
				_ = beego.AppConfig.Set("opentripmap_geoname_url", "")
			}

			resultData, err := FetchCoordinatesByCity(tc.inputCity)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an execution error, but received nil")
				}
				if resultData != nil {
					t.Errorf("Expected output payload to be nil on error states, got: %v", resultData)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no execution error, but got: %v", err)
			}

			var expectedData GeonameDTO
			if unmarshalErr := json.Unmarshal([]byte(tc.mockResponseBody), &expectedData); unmarshalErr != nil {
				t.Fatalf("Test setup failure while parsing mock JSON string: %v", unmarshalErr)
			}

			if !reflect.DeepEqual(resultData, &expectedData) {
				t.Errorf("Returned data does not match mock API exactly.\nReturned: %+v\nExpected: %+v", resultData, expectedData)
			}
		})
	}
}

func TestFetchAttractionsByCoordinates(t *testing.T) {
	testCases := []struct {
		name string
		inputLat float64
		inputLon float64
		mockStatusCode int
		mockResponseBody string
		expectError bool
		overrideURL string
		clearConfig bool
	}{
		{
			name: "Success - Valid features slice payload array returned",
			inputLat: 48.8566,
			inputLon: 2.3522,
			mockStatusCode: http.StatusOK,
			mockResponseBody: `{
				"features": [
					{
						"properties": {
							"name": "Eiffel Tower",
							"kinds": "historic,architecture"
						}
					}
				]
			}`,
			expectError: false,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "API Error - Remote server failure",
			inputLat: 48.8566,
			inputLon: 2.3522,
			mockStatusCode: http.StatusInternalServerError,
			mockResponseBody: `Internal Server Error`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "JSON Parse Error - Broken payload data structure",
			inputLat: 48.8566,
			inputLon: 2.3522,
			mockStatusCode: http.StatusOK,
			mockResponseBody: `{"features": [{broken}]}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "Connection Error - Target path is broken",
			inputLat: 48.8566,
			inputLon: 2.3522,
			mockStatusCode: http.StatusOK,
			mockResponseBody: ``,
			expectError: true,
			overrideURL: "http://invalid-url-domain-space-error.local",
			clearConfig: false,
		},
		{
			name: "Config Error - API Key completely missing",
			inputLat: 48.8566,
			inputLon: 2.3522,
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
				if !strings.Contains(r.URL.Query().Get("kinds"), "historic") {
					t.Errorf("Expected kinds parameter to contain 'historic', got: %s", r.URL.Query().Get("kinds"))
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

			_ = beego.AppConfig.Set("opentripmap_api_key", "mock_key")
			_ = beego.AppConfig.Set("opentripmap_radius_url", targetURL)

			if tc.clearConfig {
				_ = beego.AppConfig.Set("opentripmap_api_key", "")
				_ = beego.AppConfig.Set("opentripmap_radius_url", "")
			}

			resultData, err := FetchAttractionsByCoordinates(tc.inputLat, tc.inputLon)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an execution error, but received nil")
				}
				if resultData != nil {
					t.Errorf("Expected output payload to be nil on error states, got: %v", resultData)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no execution error, but got: %v", err)
			}

			var expectedData AttractionFeaturesDTO
			if unmarshalErr := json.Unmarshal([]byte(tc.mockResponseBody), &expectedData); unmarshalErr != nil {
				t.Fatalf("Test setup failure while parsing mock JSON string: %v", unmarshalErr)
			}

			if !reflect.DeepEqual(resultData, &expectedData) {
				t.Errorf("Returned data does not match mock API exactly.\nReturned: %+v\nExpected: %+v", resultData, expectedData)
			}
		})
	}
}
