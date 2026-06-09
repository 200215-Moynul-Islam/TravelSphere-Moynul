package services

import (
	"TravelSphere/data"
	"TravelSphere/models"
	"time"

	"github.com/google/uuid"
)

type WishlistService struct{}

func (s *WishlistService) AddToWishlist(username, countryName, note, status string) (models.WishlistEntry, error) {
	data.StoreMutex.Lock()
	defer data.StoreMutex.Unlock()

	entry := models.WishlistEntry{
		ID: uuid.New().String(),
		CountryName: countryName,
		Note: note,
		Status: status,
		CreatedAt: time.Now(),
	}

	data.WishlistStore[username] = append(data.WishlistStore[username], entry)
	return entry, nil
}
