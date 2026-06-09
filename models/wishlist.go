package models

import "time"

type WishlistEntry struct {
	ID string `json:"id"`
	CountryName string `json:"country_name"`
	Note string `json:"note"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
