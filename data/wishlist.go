package data

import (
	"TravelSphere/models"
	"sync"
)

var (
	WishlistStore = make(map[string][]models.WishlistEntry)
	StoreMutex sync.RWMutex
)
