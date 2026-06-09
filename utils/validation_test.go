package utils

import (
	"TravelSphere/constants"
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name string
		status string
		expected bool
	}{
		{
			name: "Valid Planned status",
			status: constants.StatusPlanned,
			expected: true,
		},
		{
			name: "Valid Visited status",
			status: constants.StatusVisited,
			expected: true,
		},
		{
			name: "Invalid status lowercase",
			status: "planned",
			expected: false,
		},
		{
			name: "Invalid status completely wrong",
			status: "UnknownStatus",
			expected: false,
		},
		{
			name: "Invalid status empty string",
			status: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsValidStatus(tt.status)
			if actual != tt.expected {
				t.Errorf("IsValidStatus(%q) = %v; want %v", tt.status, actual, tt.expected)
			}
		})
	}
}
