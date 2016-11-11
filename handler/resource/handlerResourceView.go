package resource

import (
	"net/http"

	"github.com/MISingularity/deepdash/handler"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var RoutesDataView = map[string]handler.Route{
	"Root": handler.Route{
		Name:        "Root",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: ViewIndexHandler,
	},
	"Index": handler.Route{
		Name:        "Index",
		Method:      "GET",
		Pattern:     handler.BaseURL.IndexURL,
		HandlerFunc: ViewIndexHandler,
	},
}

func ViewIndexHandler(c *gin.Context) {
	appid, _ := sessions.Default(c).Get("appid").(string)
	data := struct {
		handler.BaseViewURL
		AppID string
	}{
		handler.BaseURL,
		appid,
	}

	c.HTML(http.StatusOK, "index.html", data)
}
