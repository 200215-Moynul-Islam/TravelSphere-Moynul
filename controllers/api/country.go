package api

import (
	"TravelSphere/services"
	"net/http"
)

type CountryController struct {
    APIBaseController
}

func (c *CountryController) GetFilteredCountries() {
    searchQuery := c.GetString("search")
    regionQuery := c.GetString("region")

    service := &services.CountryService{}
    countries, err := service.GetFilteredCountries(searchQuery, regionQuery)
    if err != nil {
        c.SendError("Failed to fetch filtered country destinations", http.StatusInternalServerError)
        return
    }

    c.SendSuccess("Filtered country data list retrieved successfully", countries, http.StatusOK)
}
