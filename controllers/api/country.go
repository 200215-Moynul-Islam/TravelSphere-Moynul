package api

import (
	"TravelSphere/services"
	"net/http"
)

type CountryController struct {
    APIBaseController
}

func (c *CountryController) SearchCountriesByPartialName() {
	query := c.GetString("search")

	service := &services.CountryService{}
	countries, err := service.GetCountriesByPartialName(query)
	if err != nil {
		c.SendError("Failed to fetch search results", http.StatusInternalServerError)
		return
	}

	c.SendSuccess("Countries fetched successfully", countries, http.StatusOK)
}
