package services

import (
	"TravelSphere/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestGetCountriesByCodes(t *testing.T) {
	mockResponse := []utils.CountryDTO{
		{
			Cca3: "BGD",
			Population: 170000000,
			Region: "Asia",
		},
	}
	mockResponse[0].Name.Common = "Bangladesh"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("codes") == "ERR" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	beego.AppConfig.Set("restcountries_base_url", server.URL)

	testCases := []struct {
		name string
		codes []string
		expectError bool
		expectedLen int
	}{
		{
			name: "Empty input codes",
			codes: []string{},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "Successful API fetch and mapping",
			codes: []string{"BGD"},
			expectError: false,
			expectedLen: 1,
		},
		{
			name: "API server error response",
			codes: []string{"ERR"},
			expectError: true,
			expectedLen: 0,
		},
	}

	service := &CountryService{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestGetCountriesByPartialName(t *testing.T) {
	mockResponse := []utils.CountryDTO{
		{
			Cca3: "BGD",
			Population: 170000000,
			Region: "Asia",
		},
	}
	mockResponse[0].Name.Common = "Bangladesh"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/name/ERR") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if strings.Contains(r.URL.Path, "/name/MALFORMED") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[{invalid-json}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	_ = beego.AppConfig.Set("restcountries_base_url", server.URL)

	testCases := []struct {
		name string
		partialName string
		expectError bool
		expectedLen int
	}{
		{
			name: "Empty input search string",
			partialName: "   ",
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "Successful search match and data translation",
			partialName: "Bangla",
			expectError: false,
			expectedLen: 1,
		},
		{
			name: "API backend server internal error",
			partialName: "ERR",
			expectError: true,
			expectedLen: 0,
		},
		{
			name: "Data mapping translation failure via bad JSON structural syntax",
			partialName: "MALFORMED",
			expectError: true,
			expectedLen: 0,
		},
	}

	service := &CountryService{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.GetCountriesByPartialName(tc.partialName)

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
					t.Errorf("mapped search records inside data slice are incorrect: %+v", result[0])
				}
			}
		})
	}
}

func TestGetAllCountries(t *testing.T) {
    mockResponse := []utils.CountryDTO{
        {
            Cca3: "BGD",
            Population: 170000000,
            Region: "Asia",
        },
    }
    mockResponse[0].Name.Common = "Bangladesh"

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _ = json.NewEncoder(w).Encode(mockResponse)
    }))
    defer server.Close()

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
                    _, _ = w.Write([]byte(`[{invalid-json}`))
                })
            } else {
                server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    w.WriteHeader(http.StatusOK)
                    _ = json.NewEncoder(w).Encode(mockResponse)
                })
            }

            result, err := service.GetAllCountries()

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
