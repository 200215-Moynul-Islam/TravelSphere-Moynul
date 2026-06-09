package routers

import (
	"TravelSphere/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/countries", &controllers.CountryController{}, "get:GetAll")
	beego.Router("/countries/:code", &controllers.CountryController{}, "get:GetDetails")

	// Register global error controller
    beego.ErrorController(&controllers.ErrorController{})
}
