package routers

import (
	"TravelSphere/controllers/api"
	"TravelSphere/middlewares"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/countries",
			beego.NSRouter("", &api.CountryController{}, "get:GetFilteredCountries"),
		),
		beego.NSNamespace("/wishlist",
			beego.NSBefore(middlewares.Authenticate),
			beego.NSRouter("", &api.WishlistController{}, "get:GetWishlist;post:CreateWishlist"),
			beego.NSRouter("/:id", &api.WishlistController{}, "delete:DeleteWishlist"),
		),
	)

	beego.AddNamespace(ns)
}
