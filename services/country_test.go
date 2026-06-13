package services

import (
	"TravelSphere/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestGetCountriesByCodes(t *testing.T) {
	mockEnvelope := struct {
		Data struct {
			Objects []utils.CountryDTO `json:"objects"`
		} `json:"data"`
	}{}
	
	dto := utils.CountryDTO{
		Population: 170000000,
		Region: "Asia",
	}
	dto.Codes.Alpha3 = "BGD"
	dto.Names.Common = "Bangladesh"
	mockEnvelope.Data.Objects = append(mockEnvelope.Data.Objects, dto)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockEnvelope)
	}))
	defer server.Close()

	beego.AppConfig.Set("restcountries_api_key", "mock-key")
	beego.AppConfig.Set("restcountries_base_url", server.URL)

	testCases := []struct {
		name string
		codes []string
		expectError bool
		expectedLen int
		triggerError bool
	}{
		{
			name: "Empty input codes",
			codes: []string{},
			expectError: false,
			expectedLen: 0,
			triggerError: false,
		},
		{
			name: "Successful API fetch and mapping",
			codes: []string{"BGD"},
			expectError: false,
			expectedLen: 1,
			triggerError: false,
		},
		{
			name: "API server error response",
			codes: []string{"BGD"},
			expectError: true,
			expectedLen: 0,
			triggerError: true,
		},
	}

	service := &CountryService{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.triggerError {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})
			} else {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(mockEnvelope)
				})
			}

			result, err := service.GetCountriesByCodes(tc.codes)

			if tc.expectError {
				if err == nil {
					t.Error("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error occurred: %v", err)
			}

			if len(result) != tc.expectedLen {
				t.Errorf("expected slice length of %d, got %d", tc.expectedLen, len(result))
			}

			if tc.expectedLen > 0 {
				if result[0].Code != "BGD" || result[0].Name != "Bangladesh" {
					t.Errorf("mapped data values inside slice are incorrect: %+v", result[0])
				}
			}
		})
	}
}

func TestGetAllCountries(t *testing.T) {
	mockEnvelope := struct {
		Data struct {
			Objects []utils.CountryDTO `json:"objects"`
		} `json:"data"`
	}{}
	
	dto := utils.CountryDTO{
		Population: 170000000,
		Region: "Asia",
	}
	dto.Codes.Alpha3 = "BGD"
	dto.Names.Common = "Bangladesh"
	mockEnvelope.Data.Objects = append(mockEnvelope.Data.Objects, dto)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockEnvelope)
	}))
	defer server.Close()

	beego.AppConfig.Set("restcountries_api_key", "mock-key")
	_ = beego.AppConfig.Set("restcountries_base_url", server.URL)

	testCases := []struct {
		name string
		triggerError bool
		triggerMalformed bool
		expectError bool
		expectedLen int
	}{
		{
			name: "Successful fetch all and mapping transformation",
			triggerError: false,
			triggerMalformed: false,
			expectError: false,
			expectedLen: 1,
		},
		{
			name: "API backend server internal error fallback",
			triggerError: true,
			triggerMalformed: false,
			expectError: true,
			expectedLen: 0,
		},
		{
			name: "Translation mapping failure via bad structural syntax JSON",
			triggerError: false,
			triggerMalformed: true,
			expectError: true,
			expectedLen: 0,
		},
	}

	service := &CountryService{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.triggerError {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})
			} else if tc.triggerMalformed {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"data": {"objects": [{invalid-json}]}}`))
				})
			} else {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(mockEnvelope)
				})
			}

			result, err := service.GetAllCountries(10)

			if tc.expectError {
				if err == nil {
					t.Error("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error occurred: %v", err)
			}

			if len(result) != tc.expectedLen {
				t.Errorf("expected slice length of %d, got %d", tc.expectedLen, len(result))
			}

			if tc.expectedLen > 0 {
				if result[0].Code != "BGD" || result[0].Name != "Bangladesh" {
					t.Errorf("mapped full master records inside data slice are incorrect: %+v", result[0])
				}
			}
		})
	}
}

func TestGetFilteredCountries(t *testing.T) {
	mockEnvelopeAsia := struct {
		Data struct {
			Objects []utils.CountryDTO `json:"objects"`
		} `json:"data"`
	}{}
	dtoAsia := utils.CountryDTO{
		Population: 170000000,
		Region: "Asia",
	}
	dtoAsia.Codes.Alpha3 = "BGD"
	dtoAsia.Names.Common = "Bangladesh"
	dtoAsia.Capitals = []struct {
		Name string `json:"name"`
	}{{Name: "Dhaka"}}
	mockEnvelopeAsia.Data.Objects = append(mockEnvelopeAsia.Data.Objects, dtoAsia)

	mockEnvelopeAll := struct {
		Data struct {
			Objects []utils.CountryDTO `json:"objects"`
		} `json:"data"`
	}{}
	dtoFrance := utils.CountryDTO{
		Population: 68000000,
		Region: "Europe",
	}
	dtoFrance.Codes.Alpha3 = "FRA"
	dtoFrance.Names.Common = "France"
	dtoFrance.Capitals = []struct {
		Name string `json:"name"`
	}{{Name: "Paris"}}
	mockEnvelopeAll.Data.Objects = append(mockEnvelopeAll.Data.Objects, dtoAsia, dtoFrance)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockEnvelopeAll)
	}))
	defer server.Close()

	beego.AppConfig.Set("restcountries_api_key", "mock-key")
	_ = beego.AppConfig.Set("restcountries_base_url", server.URL)

	testCases := []struct {
		name string
		search string
		region string
		triggerError bool
		expectError bool
		expectedLen int
		expectedFirstCca string
	}{
		{
			name: "Filter by region only",
			search: "",
			region: "Asia",
			triggerError: false,
			expectError: false,
			expectedLen: 1,
			expectedFirstCca: "BGD",
		},
		{
			name: "Filter by search keyword match name",
			search: "Fran",
			region: "all",
			triggerError: false,
			expectError: false,
			expectedLen: 1,
			expectedFirstCca: "FRA",
		},
		{
			name: "Filter by search keyword match capital",
			search: "dhak",
			region: "",
			triggerError: false,
			expectError: false,
			expectedLen: 1,
			expectedFirstCca: "BGD",
		},
		{
			name: "No filters specified returns all items",
			search: "",
			region: "",
			triggerError: false,
			expectError: false,
			expectedLen: 2,
			expectedFirstCca: "BGD",
		},
		{
			name: "Filters find no matching documents",
			search: "InvalidCountryName",
			region: "Europe",
			triggerError: false,
			expectError: false,
			expectedLen: 0,
			expectedFirstCca: "",
		},
		{
			name: "Upstream backend failure scenario",
			search: "",
			region: "Asia",
			triggerError: true,
			expectError: true,
			expectedLen: 0,
			expectedFirstCca: "",
		},
	}

	service := &CountryService{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.triggerError {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				})
			} else {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					if r.URL.Query().Get("region") == "asia" {
						_ = json.NewEncoder(w).Encode(mockEnvelopeAsia)
					} else {
						_ = json.NewEncoder(w).Encode(mockEnvelopeAll)
					}
				})
			}

			result, err := service.GetFilteredCountries(tc.search, tc.region)

			if tc.expectError {
				if err == nil {
					t.Error("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error occurred: %v", err)
			}

			if len(result) != tc.expectedLen {
				t.Errorf("expected slice length of %d, got %d", tc.expectedLen, len(result))
			}

			if tc.expectedLen > 0 && tc.expectedFirstCca != "" {
				if result[0].Code != tc.expectedFirstCca {
					t.Errorf("expected first matched code to be %s, got %s", tc.expectedFirstCca, result[0].Code)
				}
			}
		})
	}
}
