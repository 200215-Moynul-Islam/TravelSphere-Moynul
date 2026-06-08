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
