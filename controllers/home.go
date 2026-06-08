package controllers

import "TravelSphere/constants"

type HomeController struct {
	SSRBaseController
}

func (c *HomeController) Get() {
	c.Data["Title"] = "Home - TravelSphere"
	c.Data["CurrentNav"] = constants.NavHome
	c.Layout = "layouts/base.tpl"
	c.TplName = "pages/home.tpl"
}
