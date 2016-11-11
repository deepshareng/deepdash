// Package handler implements a URL handler system, applying exclusive
// handler to deal with user URL request.
package client

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/MISingularity/deepdash/handler"

	"github.com/MISingularity/deepdash/pkg/log"
)

var RoutesViewNeedRedirect = map[string]handler.Route{
	"Login": handler.Route{
		Name:        "Login",
		Method:      "GET",
		Pattern:     handler.BaseURL.LoginURL,
		HandlerFunc: ViewLoginHandler,
	},
	"UnverifiedEmail": handler.Route{
		Name:        "UnverifiedEmail",
		Method:      "GET",
		Pattern:     "unverified-email",
		HandlerFunc: ViewUnverifiedEmailHandler,
	},
}

var RoutesGeneralView = map[string]handler.Route{
	"Forget": handler.Route{
		Name:        "Forget",
		Method:      "GET",
		Pattern:     handler.BaseURL.ForgetURL,
		HandlerFunc: ViewForgetHandler,
	},
	"Reset": handler.Route{
		Name:        "Reset",
		Method:      "GET",
		Pattern:     handler.BaseURL.ResetURL,
		HandlerFunc: ViewResetHandler,
	},
	"Register": handler.Route{
		Name:        "Register",
		Method:      "GET",
		Pattern:     handler.BaseURL.RegisterURL,
		HandlerFunc: ViewRegisterHandler,
	},
}

func ViewLoginHandler(c *gin.Context) {
	s := "/authorization"
	callback := c.Query("callback")
	log.Info("calllback in login is: %s\n", callback)

	data := struct {
		handler.BaseViewURL
		AuthURL          string
		CallbackURL      string
		LoginSelectorURL string
		ResendEmailURL   string
	}{
		handler.BaseURL,
		s,
		callback,
		RoutesSession["SessionLoginSelector"].Pattern,
		"/unverified-email",
	}
	c.HTML(http.StatusOK, "login.html", data)
}

func ViewUnverifiedEmailHandler(c *gin.Context) {
	session := sessions.Default(c)
	emailObj := session.Get("unverified-email")
	email := ""
	if emailObj == nil {
		email = ""
	} else {
		email = fmt.Sprint(emailObj)
	}

	data := struct {
		handler.BaseViewURL
		Email string
	}{
		handler.BaseURL,
		email,
	}
	c.HTML(http.StatusOK, "unverified-email.html", data)
}

func ViewRegisterHandler(c *gin.Context) {
	data := struct {
		handler.BaseViewURL
	}{
		handler.BaseURL,
	}
	c.HTML(http.StatusOK, "register.html", data)
}

func ViewForgetHandler(c *gin.Context) {
	data := struct {
		handler.BaseViewURL
	}{
		handler.BaseURL,
	}
	c.HTML(http.StatusOK, "forget.html", data)
}

func ViewResetHandler(c *gin.Context) {
	data := struct {
		handler.BaseViewURL
		Username string
		Token    string
	}{
		handler.BaseURL,
		c.Query("username"),
		c.Query("token"),
	}
	c.HTML(http.StatusOK, "reset.html", data)
}
