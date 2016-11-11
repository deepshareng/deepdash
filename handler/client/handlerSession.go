package client

import (
	"fmt"
	"net/http"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/handler/auth"
	"github.com/MISingularity/deepdash/pkg/errorutil"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/MISingularity/deepdash/pkg/sessionutil"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	endpoint "golang.org/x/oauth2/github"
)

var RoutesSession = map[string]handler.Route{
	"SessionLoginSelector": handler.Route{
		Name:        "SessionLogout",
		Method:      "GET",
		Pattern:     "/session/loginselector",
		HandlerFunc: SessionLoginSelectorHandler,
	},
	"SessionLogout": handler.Route{
		Name:        "SessionLogout",
		Method:      "GET",
		Pattern:     "/session/logout",
		HandlerFunc: SessionLogoutHandler,
	},
	"CheckUserName": handler.Route{
		Name:        "CheckUserName",
		Method:      "POST",
		Pattern:     "/session/checkusername",
		HandlerFunc: CheckUserNameHandler,
	},
	"GetSessionUserInfo": handler.Route{
		Name:        "GetSessionUserInfo",
		Method:      "GET",
		Pattern:     "/session/getuser",
		HandlerFunc: GetSessionUserInfoHandler,
	},
	"SetSessionAppId": handler.Route{
		Name:        "SetSessionAppId",
		Method:      "POST",
		Pattern:     "/session/setappid/:appid",
		HandlerFunc: SetSessionAppIdHandler,
	},
}

var Conf = &oauth2.Config{
	ClientID:     "9b85b1e3ef194767dcde",
	ClientSecret: "e7de8b20bdc18a0d4c221a319ef1a585b3c187a4",
	Scopes:       []string{"user:email", "repo", "openid", "profile"},
	Endpoint:     endpoint.Endpoint,
}

func SessionLoginSelectorHandler(c *gin.Context) {
	log.Infof("[Session Resource] Request to SessionLoginSelectorHandler, %s", c.Request.URL.String())
	f := httputil.NewGinframework(c)
	if auth.IsSuperAdministrator(sessionutil.GetAccountID(c)) {
		f.Redirect(http.StatusTemporaryRedirect, "/private/status")
	}
	f.Redirect(http.StatusTemporaryRedirect, "/")
}

func SessionLogoutHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	session := sessions.Default(c)
	session.Delete("Token")
	session.Delete("Username")
	session.Delete("Displayname")
	err := session.Save()
	if err != nil {
		f.Redirect(http.StatusInternalServerError, handler.BaseURL.LoginURL)
		return
	}
	f.Redirect(http.StatusTemporaryRedirect, handler.BaseURL.LoginURL)
}

func CheckUserNameHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[User Resource] Request to CheckUserNameHandler, %s", c.Request.URL.String())
	log.Infof("[User Resource] Request Detail, Username : %s", c.PostForm("username"))

	exist, err := CheckUserExist(c.PostForm("username"))

	log.Debug("CheckUser", exist, err)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	if !exist {
		f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: false})
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

func GetSessionUserInfoHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[User Resource] Request to AuthSessionUserHandler, %s", c.Request.URL.String())
	f.WriteData(http.StatusOK, struct {
		Username string `json:"username"`
	}{
		sessionutil.GetDisplayName(c),
	})

}

func SetSessionAppIdHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[User Resource] Request to SetSessionAppIdHandler, %s", c.Request.URL.String())

	appId := c.Param("appid")
	log.Infof("[User Resource] Session Request Detail, AppId : %s", appId)
	session := sessions.Default(c)
	session.Set("appid", appId)
	err := session.Save()
	if err != nil {
		f.WriteHTTPError(api.ErrSessionSave, fmt.Sprintf("[User Resource] Session save failed! Err Msg=%v", err))
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}
