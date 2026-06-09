package controllers

import (
	"TravelSphere/constants"

	beego "github.com/beego/beego/v2/server/web"
)

type SSRBaseController struct {
	beego.Controller
}

func (c *SSRBaseController) Prepare() {
	c.Data["NavItems"] = constants.NavigationItems
}
