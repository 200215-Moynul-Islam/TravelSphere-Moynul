package services

import (
	"TravelSphere/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestGetPopularAttractions(t *testing.T) {
	testCases := []struct {
		name string
		inputCity string
		mockGeonameStatus int
		mockGeonameBody string
		mockRadiusStatus int
		mockRadiusBody string
		expectError bool
		expectedResult []models.Attraction
	}{
		{
			name: "Success - Service parses coordinates and maps attractions correctly",
			inputCity: "Paris",
			mockGeonameStatus: http.StatusOK,
			mockGeonameBody: `{"lat": 48.8566, "lon": 2.3522}`,
			mockRadiusStatus: http.StatusOK,
			mockRadiusBody: `{
				"features": [
					{
						"properties": {
							"name": "Eiffel Tower",
							"kinds": "historic, architecture"
						}
					},
					{
						"properties": {
							"name": "",
							"kinds": "skipped_landmark"
						}
					}
				]
			}`,
			expectError: false,
			expectedResult: []models.Attraction{
				{
					Name: "Eiffel Tower",
					Tags: []string{"historic", "architecture"},
				},
			},
		},
		{
			name: "Edge Case - Empty city name input strings return empty slice immediately",
			inputCity: "   ",
			mockGeonameStatus: http.StatusOK,
			mockGeonameBody: ``,
			mockRadiusStatus: http.StatusOK,
			mockRadiusBody: ``,
			expectError: false,
			expectedResult: []models.Attraction{},
		},
		{
			name: "Failure - Geoname API coordinate lookup breaks",
			inputCity: "Paris",
			mockGeonameStatus: http.StatusInternalServerError,
			mockGeonameBody: `Internal Error`,
			mockRadiusStatus: http.StatusOK,
			mockRadiusBody: ``,
			expectError: true,
			expectedResult: nil,
		},
		{
			name: "Failure - Radius API landmarks lookup breaks",
			inputCity: "Paris",
			mockGeonameStatus: http.StatusOK,
			mockGeonameBody: `{"lat": 48.8566, "lon": 2.3522}`,
			mockRadiusStatus: http.StatusInternalServerError,
			mockRadiusBody: `Internal Error`,
			expectError: true,
			expectedResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = beego.AppConfig.Set("opentripmap_api_key", "mock_service_key")

			geonameServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.mockGeonameStatus)
				_, _ = w.Write([]byte(tc.mockGeonameBody))
			}))
			defer geonameServer.Close()

			radiusServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.mockRadiusStatus)
				_, _ = w.Write([]byte(tc.mockRadiusBody))
			}))
			defer radiusServer.Close()

			_ = beego.AppConfig.Set("opentripmap_geoname_url", geonameServer.URL)
			_ = beego.AppConfig.Set("opentripmap_radius_url", radiusServer.URL)

			service := &AttractionService{}
			result, err := service.GetPopularAttractions(tc.inputCity)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error execution state, but received nil")
				}
				if result != nil {
					t.Errorf("Expected result slice payload data container to be nil, got: %v", result)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected execution flow to finish without error states, got: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expectedResult) {
				t.Errorf("Returned data slice values do not match expected outcomes exactly.\nReturned: %+v\nExpected: %+v", result, tc.expectedResult)
			}
		})
	}
}
