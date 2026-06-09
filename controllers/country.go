package controllers

import (
	"TravelSphere/constants"
	"TravelSphere/models"
	"TravelSphere/services"
	"fmt"

	logs "github.com/beego/beego/v2/core/logs"
)

type CountryController struct {
	SSRBaseController
}

func (c *CountryController) GetAll() {
	c.Data["Title"] = "Explore Countries - TravelSphere"
	c.Data["CurrentNav"] = constants.NavCountries
	c.Data["PageStylesheets"] = `<link rel="stylesheet" href="/static/css/countries.css">`
	c.Data["PageScripts"] = `<script src="/static/js/countries.js"></script>`

	service := &services.CountryService{}

	allCountries, err := service.GetAllCountries()
	if err != nil {
		logs.Error("Failed to fetch country grid datasets: %v", err)
		c.Data["Countries"] = []any{}
	} else {
		c.Data["Countries"] = allCountries
	}

	c.Layout = "layouts/base.tpl"
	c.TplName = "pages/countries.tpl"
}

func (c *CountryController) GetDetails() {
	countryCode := c.Ctx.Input.Param(":code")
	if countryCode == "" {
		c.Redirect("/", 302)
		return
	}

	countryService := &services.CountryService{}
	attractionService := &services.AttractionService{}

	// Retrieve country details using the provided country code
	// Reuse the service method that fetches countries by their codes, passing a single code in the slice
	countries, err := countryService.GetCountriesByCodes([]string{countryCode})
	if err != nil || len(countries) == 0 {
		logs.Error("Failed to look up country details for code %s: %v", countryCode, err)
		c.Redirect("/", 302)
		return
	}
	country := countries[0]

	var attractions []models.Attraction
	if country.Capital != "" {
		fetchedAttractions, attrErr := attractionService.GetPopularAttractions(country.Capital)
		if attrErr != nil {
			logs.Error("Failed loading attraction cards for city %s: %v", country.Capital, attrErr)
			attractions = []models.Attraction{}
		} else {
			attractions = fetchedAttractions
		}
	}

	c.Data["Title"] = fmt.Sprintf("%s - TravelSphere", country.Name)
	c.Data["CurrentNav"] = constants.NavCountries
	c.Data["PageStylesheets"] = `<link rel="stylesheet" href="/static/css/destination.css">`

	c.Data["CountryDetails"] = country
	c.Data["Attractions"] = attractions

	c.Layout = "layouts/base.tpl"
	c.TplName = "pages/destination.tpl"
}
