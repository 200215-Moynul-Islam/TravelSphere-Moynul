package routers

import (
	"TravelSphere/controllers/api"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/countries",
			beego.NSRouter("/", &api.CountryController{}, "get:SearchCountriesByPartialName"),
		),
	)

	beego.AddNamespace(ns)
}
