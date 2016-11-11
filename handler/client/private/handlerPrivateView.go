package resource

import (
	"net/http"

	"github.com/MISingularity/deepdash/handler"
	"github.com/gin-gonic/gin"
)

var RoutesPrivateViews = map[string]handler.Route{
	"Status": handler.Route{
		Name:        "AppStatus",
		Method:      "GET",
		Pattern:     "/private/status",
		HandlerFunc: PrivateStatusViewHandler,
	},
	"VitalStatus": handler.Route{
		Name:        "VitalStatus",
		Method:      "GET",
		Pattern:     "/private/vital",
		HandlerFunc: PrivateVitalStatusViewHandler,
	},
	"AddUser": handler.Route{
		Name:        "AddUser",
		Method:      "GET",
		Pattern:     "/private/add-user",
		HandlerFunc: PrivateUserAddViewHandler,
	},
	"Permission": handler.Route{
		Name:        "Permission",
		Method:      "GET",
		Pattern:     "/private/permission",
		HandlerFunc: PermissionViewHandler,
	},
}

type PrivateView struct {
	StatusView string
}

var BasicPrivateView = PrivateView{
	StatusView: "/private/status",
}

func PermissionViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "private-permission.html", BasicPrivateView)
}

func PrivateUserAddViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "private-add-user.html", BasicPrivateView)
}

func PrivateStatusViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "private-status.html", BasicPrivateView)
}

func PrivateVitalStatusViewHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "private-vital-status.html", BasicPrivateView)
}
