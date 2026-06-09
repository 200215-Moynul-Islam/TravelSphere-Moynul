package api

import beego "github.com/beego/beego/v2/server/web"

type APIBaseController struct {
	beego.Controller
}

type ApiResponse struct {
	Success bool `json:"success"`
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func (c *APIBaseController) SendSuccess(message string, data any, statusCode int) {
	c.Data["json"] = ApiResponse{
		Success: true,
		Message: message,
		Data: data,
	}
	c.Ctx.Output.SetStatus(statusCode)
	c.ServeJSON()
}

func (c *APIBaseController) SendError(message string, statusCode int) {
	c.Data["json"] = ApiResponse{
		Success: false,
		Error: message,
	}
	c.Ctx.Output.SetStatus(statusCode)
	c.ServeJSON()
}
