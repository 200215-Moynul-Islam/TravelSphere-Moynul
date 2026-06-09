package services

import (
	"TravelSphere/data"
	"TravelSphere/models"
	"testing"
	"time"
)

func TestAddToWishlist(t *testing.T) {
	tests := []struct {
		name string
		username string
		countryName string
		note string
		status string
	}{
		{
			name: "Valid Planned Destination",
			username: "moynul_islam",
			countryName: "Japan",
			note: "Cherry blossom season",
			status: "Planned",
		},
		{
			name: "Valid Visited Destination",
			username: "moynul_islam",
			countryName: "France",
			note: "Eiffel tower trip",
			status: "Visited",
		},
		{
			name: "Empty Note Field",
			username: "john_doe",
			countryName: "Canada",
			note: "",
			status: "Planned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.StoreMutex.Lock()
			data.WishlistStore = make(map[string][]models.WishlistEntry)
			data.StoreMutex.Unlock()

			service := &WishlistService{}
			entry, err := service.AddToWishlist(tt.username, tt.countryName, tt.note, tt.status)

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
			userEntries, exists := data.WishlistStore[tt.username]
			data.StoreMutex.RUnlock()

			if !exists {
				t.Fatalf("Expected entries slice to exist for username %q", tt.username)
			}

			found := false
			for _, storedEntry := range userEntries {
				if storedEntry.ID == entry.ID {
					found = true
					if storedEntry.CountryName != tt.countryName {
						t.Errorf("Stored entry mismatch: expected country %q, got %q", tt.countryName, storedEntry.CountryName)
					}
					if storedEntry.Note != tt.note {
						t.Errorf("Stored entry mismatch: expected note %q, got %q", tt.note, storedEntry.Note)
					}
					if storedEntry.Status != tt.status {
						t.Errorf("Stored entry mismatch: expected status %q, got %q", tt.status, storedEntry.Status)
					}
					break
				}
			}

			if !found {
				t.Errorf("Expected entry with ID %q to exist inside data.WishlistStore slice for user %q", entry.ID, tt.username)
			}
		})
	}
}

func TestDeleteWishlist(t *testing.T) {
	tests := []struct {
		name string
		username string
		targetID string
		initialStore map[string][]models.WishlistEntry
		expectedErr bool
		expectedLen int
	}{
		{
			name: "Successful deletion",
			username: "moynul_islam",
			targetID: "id-123",
			initialStore: map[string][]models.WishlistEntry{
				"moynul_islam": {
					{ID: "id-123", CountryName: "Japan", CreatedAt: time.Now()},
					{ID: "id-456", CountryName: "France", CreatedAt: time.Now()},
				},
			},
			expectedErr: false,
			expectedLen: 1,
		},
		{
			name: "Entry ID not found",
			username: "moynul_islam",
			targetID: "id-unknown",
			initialStore: map[string][]models.WishlistEntry{
				"moynul_islam": {
					{ID: "id-123", CountryName: "Japan", CreatedAt: time.Now()},
				},
			},
			expectedErr: true,
			expectedLen: 1,
		},
		{
			name: "User has no entries",
			username: "unknown_user",
			targetID: "id-123",
			initialStore: map[string][]models.WishlistEntry{
				"moynul_islam": {
					{ID: "id-123", CountryName: "Japan", CreatedAt: time.Now()},
				},
			},
			expectedErr: true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.StoreMutex.Lock()
			data.WishlistStore = tt.initialStore
			data.StoreMutex.Unlock()

			service := &WishlistService{}
			err := service.DeleteWishlist(tt.username, tt.targetID)

			if (err != nil) != tt.expectedErr {
				t.Fatalf("Expected error presence: %v, got: %v", tt.expectedErr, err)
			}

			data.StoreMutex.RLock()
			remainingLen := len(data.WishlistStore[tt.username])
			data.StoreMutex.RUnlock()

			if !tt.expectedErr && remainingLen != tt.expectedLen {
				t.Errorf("Expected remaining slice length to be %d, got %d", tt.expectedLen, remainingLen)
			}
		})
	}
}

func TestGetWishlist(t *testing.T) {
	tests := []struct {
		name string
		username string
		initialStore map[string][]models.WishlistEntry
		expectedLen int
	}{
		{
			name: "Retrieve multiple entries for user",
			username: "moynul_islam",
			initialStore: map[string][]models.WishlistEntry{
				"moynul_islam": {
					{ID: "id-1", CountryName: "Japan", CreatedAt: time.Now()},
					{ID: "id-2", CountryName: "France", CreatedAt: time.Now()},
				},
			},
			expectedLen: 2,
		},
		{
			name: "Retrieve empty slice for non-existent user",
			username: "new_user",
			initialStore: map[string][]models.WishlistEntry{
				"moynul_islam": {
					{ID: "id-1", CountryName: "Japan", CreatedAt: time.Now()},
				},
			},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.StoreMutex.Lock()
			data.WishlistStore = tt.initialStore
			data.StoreMutex.Unlock()

			service := &WishlistService{}
			entries, err := service.GetWishlist(tt.username)

			if err != nil {
				t.Fatalf("Expected no error, but got: %v", err)
			}
			if len(entries) != tt.expectedLen {
				t.Errorf("Expected entries length to be %d, got %d", tt.expectedLen, len(entries))
			}
			if len(entries) == 0 && entries == nil {
				t.Error("Expected an allocated empty slice, but got nil")
			}
		})
	}
}

func TestUpdateWishlist(t *testing.T) {
	tests := []struct {
		name string
		username string
		targetID string
		inputNote string
		inputStatus string
		initialStore map[string][]models.WishlistEntry
		expectedErr bool
	}{
		{
			name: "Successful update modification",
			username: "moynul_islam",
			targetID: "id-123",
			inputNote: "Updated note text",
			inputStatus: "Visited",
			initialStore: map[string][]models.WishlistEntry{
				"moynul_islam": {
					{ID: "id-123", CountryName: "Japan", Note: "Old note", Status: "Planned", CreatedAt: time.Now()},
				},
			},
			expectedErr: false,
		},
		{
			name: "Target item execution failure missing ID",
			username: "moynul_islam",
			targetID: "id-missing",
			inputNote: "Target note",
			inputStatus: "Visited",
			initialStore: map[string][]models.WishlistEntry{
				"moynul_islam": {
					{ID: "id-123", CountryName: "Japan", Note: "Old note", Status: "Planned", CreatedAt: time.Now()},
				},
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.StoreMutex.Lock()
			data.WishlistStore = tt.initialStore
			data.StoreMutex.Unlock()

			service := &WishlistService{}
			entry, err := service.UpdateWishlist(tt.username, tt.targetID, tt.inputNote, tt.inputStatus)

			if (err != nil) != tt.expectedErr {
				t.Fatalf("Expected error presence: %v, got: %v", tt.expectedErr, err)
			}

			if !tt.expectedErr {
				if entry.Note != tt.inputNote {
					t.Errorf("Expected note to be updated to %q, got %q", tt.inputNote, entry.Note)
				}
				if entry.Status != tt.inputStatus {
					t.Errorf("Expected status to be updated to %q, got %q", tt.inputStatus, entry.Status)
				}

				data.StoreMutex.RLock()
				storedEntries := data.WishlistStore[tt.username]
				data.StoreMutex.RUnlock()

				if storedEntries[0].Note != tt.inputNote || storedEntries[0].Status != tt.inputStatus {
					t.Error("Stored entry inside data layer slice map was not structurally mutated")
				}
			}
		})
	}
}
