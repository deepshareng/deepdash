package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/config"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/sessionutil"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
	"gopkg.in/mgo.v2/bson"
)

func extractToken(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("Invalid authentication")
	}

	return parts[1], nil
}

func IsSuperAdministrator(accountid string) bool {
	if !bson.IsObjectIdHex(accountid) {
		return false
	}
	st, err := storage.MongoUS.GetUser(storage.UserFormat{Id: bson.ObjectIdHex(accountid)}, false, []string{}, []string{})
	if err != nil {
		return false
	}
	permissionList := strings.Split(st[0].Permission, "_")
	for i := 0; i < len(permissionList); i++ {
		if permissionList[i] == "administrator" {
			return true
		}
	}
	return false
}

func IsDemoAccount(accountid string) bool {
	if !bson.IsObjectIdHex(accountid) {
		return false
	}
	st, err := storage.MongoUS.GetUser(storage.UserFormat{Id: bson.ObjectIdHex(accountid)}, false, []string{}, []string{})
	if err != nil {
		return false
	}
	permissionList := strings.Split(st[0].Permission, "_")
	for i := 0; i < len(permissionList); i++ {
		if permissionList[i] == "demo" {
			return true
		}
	}
	return false
}

func TokenAuthentication(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	v := session.Get("Token")
	if v == nil {
		return "No Authorization token!", false
	}
	accountid := sessionutil.GetAccountID(c)
	if c.Request.Method != "GET" && IsSuperAdministrator(accountid) {
		return "", false
	}
	if c.Request.Method != "GET" && IsDemoAccount(accountid) {
		return "", false
	}
	appid := c.Param("appid")
	if IsSuperAdministrator(sessionutil.GetAccountID(c)) {
		appid = ""
	}

	if IsDemoAccount(sessionutil.GetAccountID(c)) {
		// demo account has auth of 'linux command': which appid is: 38CCA4C77072DDC9
		// appid == ""  when login and retrive html
		if appid == "38CCA4C77072DDC9" || appid == "" {
			return "", true
		} else {
			return "", false
		}
	}

	res, err := http.Get(
		config.Cliconfig.TokenUrl +
			fmt.Sprintf("?token=%s&accountid=%s&appid=%s",
				url.QueryEscape(fmt.Sprint(v)),
				url.QueryEscape(accountid),
				url.QueryEscape(appid),
			),
	)
	if err != nil || res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		log.Info(c.Request.URL.String() + "  authentication fail!")
		if err != nil {
			log.Info("error:", err)
		} else {
			log.Info("status code:", res.StatusCode)
		}
		log.Info(config.Cliconfig)
		return "Token authentication fail!", false
	}
	return "", true
}

func ResourceViewAuthentication(c *gin.Context) {
	f := httputil.NewGinframework(c)
	//callback := c.Request.URL.String()
	callback := c.Request.URL.String()

	log.Info("callback:", url.QueryEscape(callback))
	_, ok := TokenAuthentication(c)
	if !ok {
		log.Info("callback:", url.QueryEscape(callback))
		fmt.Printf("callbak in fmt: %s\n", url.QueryEscape(callback))
		f.Redirect(
			http.StatusTemporaryRedirect,
			handler.BaseURL.LoginURL+"?callback="+url.QueryEscape(callback))
	}
	return
}

func ResourceAuthentication(c *gin.Context) {
	f := httputil.NewGinframework(c)
	errMsg, ok := TokenAuthentication(c)
	if !ok {
		f.WriteHTTPError(api.ErrAuthResourceFail, errMsg)
		c.Abort()
	}
	return
}

func GeneralViewAuthentication(c *gin.Context) {
	f := httputil.NewGinframework(c)
	callback := c.Query("callback")
	if callback == "" {
		callback = handler.BaseURL.IndexURL
		if IsSuperAdministrator(sessionutil.GetAccountID(c)) {
			callback = "/private/status"
		}
	}

	_, ok := TokenAuthentication(c)
	if ok {
		f.Redirect(http.StatusTemporaryRedirect, callback)
	}
	return
}

func AdministratorTokenAuthentication(c *gin.Context) (string, bool) {
	session := sessions.Default(c)
	accountid := sessionutil.GetAccountID(c)
	appid := c.Param("appid")
	v := session.Get("Token")
	if v == nil {
		return "No Authorization token!", false
	}
	if IsSuperAdministrator(sessionutil.GetAccountID(c)) {
		appid = ""
	}

	res, err := http.Get(
		config.Cliconfig.TokenUrl +
			fmt.Sprintf("?token=%s&accountid=%s&appid=%s",
				url.QueryEscape(fmt.Sprint(v)),
				url.QueryEscape(accountid),
				url.QueryEscape(appid),
			),
	)
	if err != nil || res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		log.Info(c.Request.URL.String() + " authentication fail!")
		return "Token authentication fail!", false
	}
	return "", true
}

func AdministratorResourceAuthentication(c *gin.Context) {
	f := httputil.NewGinframework(c)
	_, ok := AdministratorTokenAuthentication(c)
	if !ok || !IsSuperAdministrator(sessionutil.GetAccountID(c)) {
		f.WriteHTTPError(api.ErrAuthResourceFail, "权限不足")
		c.Abort()
	}
	return
}

func AdministratorViewAuthentication(c *gin.Context) {
	f := httputil.NewGinframework(c)
	callback := c.Request.URL.String()

	_, ok := AdministratorTokenAuthentication(c)
	if !ok || !IsSuperAdministrator(sessionutil.GetAccountID(c)) {
		f.Redirect(
			http.StatusTemporaryRedirect,
			handler.BaseURL.LoginURL+"?callback="+url.QueryEscape(callback))
	}
	return
}
