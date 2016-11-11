package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/config"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/RangelReale/osin"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var AuthURL string

func DownloadAccessToken(url string, auth *osin.BasicAuth, output map[string]interface{}) error {
	// download access token
	log.Info("[User Agent] DownloadAccessToken", url, auth, output)
	preq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if auth != nil {
		preq.SetBasicAuth(auth.Username, auth.Password)
	}

	pclient := &http.Client{}
	presp, err := pclient.Do(preq)
	if err != nil {
		return err
	}

	if presp.StatusCode != 200 {
		return errors.New(fmt.Sprintln("Invalid status code=", presp.StatusCode))
	}
	jdec := json.NewDecoder(presp.Body)
	err = jdec.Decode(&output)
	return err
}

func HandlerAuthRequest(c *gin.Context) {

	log.Infof("[Client] Request to HandlerAuthRequest, %s", c.Request.URL.String())

	f := httputil.NewGinframework(c)

	oauthURL := fmt.Sprintf(
		"%s/auth/authorize?response_type=code&client_id=%s&state=xyz&scope=personal&redirect_uri=%s&username=%s&password=%s",
		AuthURL,
		url.QueryEscape(config.Cliconfig.ClientId),
		url.QueryEscape(config.Cliconfig.RedirectUrl),
		url.QueryEscape(c.PostForm("username")),
		url.QueryEscape(c.PostForm("password")),
	)

	log.Debugf("[Client] Request to %s", oauthURL)

	res, err := http.Get(oauthURL)

	log.Debug("[Client] ", res, err)

	if err != nil || res.StatusCode != http.StatusOK {

		// check wheather is unactive user
		notActive, user, rightPassword, err := CheckUserIsNotActive(c.PostForm("username"), c.PostForm("password"))
		if err == nil {
			if notActive == true && rightPassword {
				session := sessions.Default(c)
				session.Set("unverified-email", user.Username)
				err = session.Save()
				if err != nil {
					log.Infof("[Client] Save Session Failed, err=%v", err)
					f.WriteHTTPError(api.ErrAuthLoginFail, "")
					return
				}
				f.WriteData(http.StatusOK, api.ResponseData{Success: false, Data: "unverified_user"})
				return
			}
		}
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Infof("[Client] Read response data Failed, err=%v", err)
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}
	auths := api.AuthSuccessValues{}
	err = json.Unmarshal(resBody, &auths)
	if err != nil {
		log.Infof("[Client] Convert response data Failed, err=%v", err)
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}
	session := sessions.Default(c)
	session.Set("Token", auths.Token)
	session.Set("Accountid", auths.Username)
	session.Set("Displayname", auths.Displayname)
	session.Delete("unverified-email")
	err = session.Save()
	if err != nil {
		log.Infof("[Client] Save Session Failed, err=%v", err)
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

// This user agent is for deepdash application,
// as a result, we use our handler redirect url to request token
// As same as the deepdash, other application will render their
// redirect url in their body when they request token from authenticaiton service
func HandlerRedirectURIRetrieveToken(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Client] Request to HandlerRedirectURIRetrieveToken, %s", c.Request.URL.String())
	code := c.Query("code")
	res, err := http.Get(fmt.Sprintf("%s/auth/getuserdata?type=code&key=%s", AuthURL, code))
	if err != nil || res.StatusCode != http.StatusOK {
		log.Infof("[Client] Load Access Failed, statuscode = %v, err=%v", res.StatusCode, err)
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}
	userdata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Infof("[Client] Convert User Data Failed, err=%v", err)
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}
	resdata := api.OauthUserData{}
	err = json.Unmarshal(userdata, &resdata)
	if err != nil {
		log.Infof("[Client] Convert User Data Failed, err=%v", err)
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}

	log.Info("[Client] Userdata=", string(userdata))
	username := resdata.AccountID
	displayname := resdata.Username
	aurl := fmt.Sprintf(
		"%s?grant_type=authorization_code&client_id=%s&client_secret=%s&state=xyz&redirect_uri=%s&code=%s",
		config.Cliconfig.AuthorizeUrl,
		url.QueryEscape(config.Cliconfig.ClientId),
		url.QueryEscape(config.Cliconfig.ClientSecret),
		url.QueryEscape(config.Cliconfig.RedirectUrl),
		url.QueryEscape(code),
	)
	jr := make(map[string]interface{})
	err = DownloadAccessToken(aurl, &osin.BasicAuth{Username: config.Cliconfig.ClientId, Password: config.Cliconfig.ClientSecret}, jr)
	if err != nil {
		log.Infof("[Client] Download Access Token Failed, err=%v", err)
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}
	if _, ok := jr["error"]; ok {
		log.Infof("[Client] Load Access Failed, err=%v", jr["error"])
		f.WriteHTTPError(api.ErrAuthLoginFail, "")
		return
	}

	token := ""
	if at, ok := jr["access_token"]; ok {
		token = fmt.Sprintf("%v", at)
	}

	f.WriteData(
		http.StatusOK,
		api.AuthSuccessValues{
			Token:       token,
			Username:    username,
			Displayname: displayname,
		},
	)
}
