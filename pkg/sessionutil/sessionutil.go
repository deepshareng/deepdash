package sessionutil

import (
	"reflect"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetDisplayName(c *gin.Context) string {
	return reflect.ValueOf(sessions.Default(c).Get("Displayname")).String()
}

func GetAccountID(c *gin.Context) string {
	return reflect.ValueOf(sessions.Default(c).Get("Accountid")).String()
}
