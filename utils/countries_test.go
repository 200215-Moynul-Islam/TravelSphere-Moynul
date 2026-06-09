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

func TestGetCountriesByCodes(t *testing.T) {
	testCases := []struct {
		name string
		inputCodes []string
		mockStatusCode int
		mockResponseBody string
		expectError bool
		overrideURL string
		clearConfig bool
	}{
		{
			name: "Success - Valid country array returned",
			inputCodes: []string{"usa", "fra", "jpn"},
			mockStatusCode: http.StatusOK,
			mockResponseBody: `[
				{"name": {"common": "United States"}, "cca3": "USA", "capital": ["Washington, D.C."]},
				{"name": {"common": "France"}, "cca3": "FRA", "capital": ["Paris"]}
			]`,
			expectError: false,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "API Error - Remote server returns 500 status",
			inputCodes: []string{"usa"},
			mockStatusCode: http.StatusInternalServerError,
			mockResponseBody: `{"message": "Internal Server Error"}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "JSON Parse Error - Server returns invalid JSON syntax",
			inputCodes: []string{"usa"},
			mockStatusCode: http.StatusOK,
			mockResponseBody: `[{invalid-json}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "Connection Error - Server URL protocol is completely invalid",
			inputCodes: []string{"usa"},
			mockStatusCode: http.StatusOK,
			mockResponseBody: ``,
			expectError: true,
			overrideURL: "http://invalid-url-domain-space-error.local",
			clearConfig: false,
		},
		{
			name: "Config Error Fallback - Triggers config missing branch and breaks on live fallback connection",
			inputCodes: []string{"invalid-country-code-payload-force-failure"},
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
				expectedJoined := strings.Join(tc.inputCodes, ",")
				if !strings.Contains(r.URL.RawQuery, "codes="+expectedJoined) {
					t.Errorf("Expected URL query parameter to contain codes=%s, but got: %s", expectedJoined, r.URL.RawQuery)
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
			_ = beego.AppConfig.Set("restcountries_base_url", targetURL)

			if tc.clearConfig {
				_ = beego.AppConfig.Set("restcountries_base_url", "")
			}

			resultData, err := GetCountriesByCodes(tc.inputCodes)

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

			var expectedData []CountryDTO
			if unmarshalErr := json.Unmarshal([]byte(tc.mockResponseBody), &expectedData); unmarshalErr != nil {
				t.Fatalf("Test setup failure while parsing mock JSON text string: %v", unmarshalErr)
			}

			if !reflect.DeepEqual(resultData, expectedData) {
				t.Errorf("Returned data does not match the mock API data exactly.\nReturned: %+v\nExpected: %+v", resultData, expectedData)
			}
		})
	}
}

func TestGetCountriesByPartialName(t *testing.T) {
	testCases := []struct {
		name string
		inputSearch string
		mockStatusCode int
		mockResponseBody string
		expectError bool
		overrideURL string
		clearConfig bool
	}{
		{
			name: "Success - Valid countries matching partial name returned",
			inputSearch: "united",
			mockStatusCode: http.StatusOK,
			mockResponseBody: `[
				{"name": {"common": "United States"}, "cca3": "USA", "capital": ["Washington, D.C."]},
				{"name": {"common": "United Kingdom"}, "cca3": "GBR", "capital": ["London"]}
			]`,
			expectError: false,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "API Error - Remote server returns 500 status on search fault",
			inputSearch: "error-trigger",
			mockStatusCode: http.StatusInternalServerError,
			mockResponseBody: `{"message": "Internal Server Error"}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "JSON Parse Error - Search endpoint returns broken syntax body",
			inputSearch: "bad-json",
			mockStatusCode: http.StatusOK,
			mockResponseBody: `[{invalid-json}`,
			expectError: true,
			overrideURL: "",
			clearConfig: false,
		},
		{
			name: "Connection Error - Server URL network target path is broken",
			inputSearch: "network-fail",
			mockStatusCode: http.StatusOK,
			mockResponseBody: ``,
			expectError: true,
			overrideURL: "http://invalid-url-domain-space-error.local",
			clearConfig: false,
		},
		{
			name: "Config Error Fallback - Missing configuration properties break client execution",
			inputSearch: "config-break",
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
				if !strings.Contains(r.URL.Path, "/name/"+tc.inputSearch) {
					t.Errorf("Expected URL path to contain /name/%s, but got: %s", tc.inputSearch, r.URL.Path)
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
			_ = beego.AppConfig.Set("restcountries_base_url", targetURL)

			if tc.clearConfig {
				_ = beego.AppConfig.Set("restcountries_base_url", "")
			}

			resultData, err := GetCountriesByPartialName(tc.inputSearch)

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

			var expectedData []CountryDTO
			if unmarshalErr := json.Unmarshal([]byte(tc.mockResponseBody), &expectedData); unmarshalErr != nil {
				t.Fatalf("Test setup failure while parsing mock JSON text string: %v", unmarshalErr)
			}

			if !reflect.DeepEqual(resultData, expectedData) {
				t.Errorf("Returned data does not match the mock API data exactly.\nReturned: %+v\nExpected: %+v", resultData, expectedData)
			}
		})
	}
}

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
            mockResponseBody: `[
                {"name": {"common": "Albania"}, "cca3": "ALB", "capital": ["Tirana"]},
                {"name": {"common": "Algeria"}, "cca3": "DStandard", "capital": ["Algiers"]}
            ]`,
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
            mockResponseBody: `[{invalid-json}`,
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
            overrideURL: "",
            clearConfig: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if !strings.HasSuffix(r.URL.Path, "/all") {
                    t.Errorf("Expected URL path to end with /all, but got: %s", r.URL.Path)
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
            _ = beego.AppConfig.Set("restcountries_base_url", targetURL)

            if tc.clearConfig {
                _ = beego.AppConfig.Set("restcountries_base_url", "")
            }

            resultData, err := GetAllCountries()

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

            var expectedData []CountryDTO
            if unmarshalErr := json.Unmarshal([]byte(tc.mockResponseBody), &expectedData); unmarshalErr != nil {
                t.Fatalf("Test setup failure while parsing mock JSON text string: %v", unmarshalErr)
            }

            if !reflect.DeepEqual(resultData, expectedData) {
                t.Errorf("Returned data does not match the mock API data exactly.\nReturned: %+v\nExpected: %+v", resultData, expectedData)
            }
        })
    }
}
