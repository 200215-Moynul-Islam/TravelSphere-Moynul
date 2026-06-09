package services

import (
	"TravelSphere/data"
	"TravelSphere/models"
	"errors"
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

func (s *WishlistService) DeleteWishlist(username, id string) error {
	data.StoreMutex.Lock()
	defer data.StoreMutex.Unlock()
	userEntries, exists := data.WishlistStore[username]
	if !exists {
		return errors.New("wishlist entry not found")
	}
	for i, entry := range userEntries {
		if entry.ID == id {
			data.WishlistStore[username] = append(userEntries[:i], userEntries[i+1:]...)
			return nil
		}
	}
	return errors.New("wishlist entry not found")
}

func (s *WishlistService) GetWishlist(username string) ([]models.WishlistEntry, error) {
	data.StoreMutex.RLock()
	defer data.StoreMutex.RUnlock()
	entries, exists := data.WishlistStore[username]
	if !exists {
		return []models.WishlistEntry{}, nil
	}
	return entries, nil
}

func (s *WishlistService) UpdateWishlist(username, id, note, status string) (models.WishlistEntry, error) {
	data.StoreMutex.Lock()
	defer data.StoreMutex.Unlock()
	userEntries, exists := data.WishlistStore[username]
	if !exists {
		return models.WishlistEntry{}, errors.New("wishlist entry not found")
	}
	for i, entry := range userEntries {
		if entry.ID == id {
			userEntries[i].Note = note
			userEntries[i].Status = status
			return userEntries[i], nil
		}
	}
	return models.WishlistEntry{}, errors.New("wishlist entry not found")
}
