package controllers

type ErrorController struct {
    SSRBaseController
}

func (c *ErrorController) Error404() {
    c.Data["Title"] = "404 - Page Not Found"
    c.Layout = "layouts/error.tpl"
    c.TplName = "pages/404.tpl"
}

func (c *ErrorController) Error500() {
    c.Data["Title"] = "500 - Internal Server Error"
    c.Layout = "layouts/error.tpl"
    c.TplName = "pages/500.tpl"
}
