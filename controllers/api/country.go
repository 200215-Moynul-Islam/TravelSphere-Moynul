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
        c.Data["json"] = map[string]string{"error": "Failed to fetch search results"}
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.ServeJSON()
        return
    }

    c.Data["json"] = countries
    c.ServeJSON()
}