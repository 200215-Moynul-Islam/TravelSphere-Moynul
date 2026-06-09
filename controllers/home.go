package controllers

import (
	"TravelSphere/constants"
	"TravelSphere/services"

	logs "github.com/beego/beego/v2/core/logs"
)

type HomeController struct {
	SSRBaseController
}

func (c *HomeController) Get() {
	c.Data["Title"] = "Home - TravelSphere"
	c.Data["CurrentNav"] = constants.NavHome
	c.Data["PageStylesheets"] = `<link rel="stylesheet" href="/static/css/home.css">`
	c.Data["PageScripts"] = `<script src="/static/js/home.js"></script>`

	service := &services.CountryService{}

	featuredCountries, err := service.GetCountriesByCodes(constants.FeaturedCountryCodes)
	if err != nil {
		logs.Error("Failed to load featured countries:", err)
		c.Data["FeaturedCountries"] = []any{}
	} else {
		c.Data["FeaturedCountries"] = featuredCountries
	}
	
	c.Layout = "layouts/base.tpl"
	c.TplName = "pages/home.tpl"
}
