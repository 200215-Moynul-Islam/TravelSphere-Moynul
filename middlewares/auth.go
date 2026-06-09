package middlewares

import (
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

func Authenticate(ctx *context.Context) {
	username := ctx.Input.Header("Username")
	username = strings.TrimSpace(username)
	if username == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]any{
			"success": false,
			"message": "User is not authenticated",
		}, false, false)
		return
	}
	ctx.Request.Header.Set("Username", username)
}
