package utils

import (
	"testing"
)

// func TestMapToCountryModel(t *testing.T) {
// 	dto := CountryDTO{
// 		Cca3: "BGD",
// 		Population: 170000000,
// 		Region: "Asia",
// 		Subregion: "Southern Asia",
// 		Capital: []string{"Dhaka"},
// 		Currencies: map[string]struct{ Name string `json:"name"` }{
// 			"BDT": {Name: "Bangladeshi taka"},
// 		},
// 		Languages: map[string]string{
// 			"ben": "Bengali",
// 		},
// 	}
// 	dto.Name.Common = "Bangladesh"
// 	dto.Name.Official = "People's Republic of Bangladesh"
// 	dto.Flags.Png = "https://flagcdn.com/w320/bd.png"

// 	expected := models.Country{
// 		Code: "BGD",
// 		Name: "Bangladesh",
// 		OfficialName: "People's Republic of Bangladesh",
// 		Flag: "https://flagcdn.com/w320/bd.png",
// 		Capital: "Dhaka",
// 		Population: "170.0M",
// 		Region: "Asia - Southern Asia",
// 		Currency: "BDT (Bangladeshi taka)",
// 		Languages: "Bengali",
// 	}

// 	result := MapToCountryModel(dto)

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("MapToCountryModel returned unexpected structure.\nGot: %+v\nExpected: %+v", result, expected)
// 	}
// }

func TestMapToCountrySlice(t *testing.T) {
	dtos := []CountryDTO{
		{
			Cca3: "USA",
			Population: 33100000,
			Region: "Americas",
		},
	}
	dtos[0].Name.Common = "United States"

	result, err := MapToCountrySlice(dtos)
	if err != nil {
		t.Fatalf("MapToCountrySlice returned an unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected slice length of 1, got %d", len(result))
	}

	if result[0].Code != "USA" || result[0].Name != "United States" {
		t.Errorf("Mapped data inside slice is incorrect: %+v", result[0])
	}
}

func TestFormatPopulation(t *testing.T) {
	testCases := []struct {
		name string
		input int64
		expected string
	}{
		{
			name: "Millions formatting",
			input: 2500000,
			expected: "2.5M",
		},
		{
			name: "Thousands formatting",
			input: 85400,
			expected: "85.4K",
		},
		{
			name: "Raw integer formatting under thousand",
			input: 420,
			expected: "420",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formatPopulation(tc.input)
			if result != tc.expected {
				t.Errorf("formatPopulation(%d) = %s; expected %s", tc.input, result, tc.expected)
			}
		})
	}
}
