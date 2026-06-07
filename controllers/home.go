package controllers

type HomeController struct {
	SSRBaseController
}

func (c *HomeController) Get() {
	c.Data["Title"] = "Home - TravelSphere"
	c.Layout = "layouts/base.tpl"
	c.TplName = "pages/home.tpl"
}
