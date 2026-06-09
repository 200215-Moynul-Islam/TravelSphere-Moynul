package services

import (
	"TravelSphere/data"
	"TravelSphere/models"
	"testing"
)

func TestAddToWishlist_TableDriven(t *testing.T) {
	tests := []struct {
		name string
		countryName string
		note string
		status string
	}{
		{
			name: "Valid Planned Destination",
			countryName: "Japan",
			note: "Cherry blossom season",
			status: "Planned",
		},
		{
			name: "Valid Visited Destination",
			countryName: "France",
			note: "Eiffel tower trip",
			status: "Visited",
		},
		{
			name: "Empty Note Field",
			countryName: "Canada",
			note: "",
			status: "Planned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.StoreMutex.Lock()
			data.WishlistStore = make(map[string]models.WishlistEntry)
			data.StoreMutex.Unlock()

			service := &WishlistService{}
			entry, err := service.AddToWishlist(tt.countryName, tt.note, tt.status)

			if err != nil {
				t.Fatalf("Expected no error, but got: %v", err)
			}
			if entry.ID == "" {
				t.Error("Expected a generated UUID ID, but got an empty string")
			}
			if entry.CountryName != tt.countryName {
				t.Errorf("Expected CountryName %q, got %q", tt.countryName, entry.CountryName)
			}
			if entry.Note != tt.note {
				t.Errorf("Expected Note %q, got %q", tt.note, entry.Note)
			}
			if entry.Status != tt.status {
				t.Errorf("Expected Status %q, got %q", tt.status, entry.Status)
			}
			if entry.CreatedAt.IsZero() {
				t.Error("Expected CreatedAt timestamp to be initialized, but got zero time")
			}

			data.StoreMutex.RLock()
			storedEntry, exists := data.WishlistStore[entry.ID]
			data.StoreMutex.RUnlock()

			if !exists {
				t.Fatalf("Expected entry with ID %q to exist in data.WishlistStore", entry.ID)
			}
			if storedEntry.CountryName != tt.countryName {
				t.Errorf("Stored entry mismatch: expected country %q, got %q", tt.countryName, storedEntry.CountryName)
			}
		})
	}
}
