package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"TravelSphere/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{})

	// Register global error controller
    beego.ErrorController(&controllers.ErrorController{})
}
