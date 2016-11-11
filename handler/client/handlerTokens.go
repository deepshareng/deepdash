package client

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/email"
	"github.com/MISingularity/deepdash/pkg/errorutil"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/MISingularity/deepdash/pkg/validator"
	"github.com/MISingularity/deepdash/storage"
	"github.com/MISingularity/deepshare2/pkg/path"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var RoutesToken = map[string]handler.Route{
	"ForgotPassword": handler.Route{
		Name:        "ForgotPassword",
		Method:      "POST",
		Pattern:     "/password/forgotpassword",
		HandlerFunc: ForgotPasswordHandler,
	},
	"ResetPassword": handler.Route{
		Name:        "ResetPassword",
		Method:      "POST",
		Pattern:     "/password/resetpassword",
		HandlerFunc: ResetPasswordHandler,
	},
	"ActivateAccount": handler.Route{
		Name:        "ActivateAccount",
		Method:      "GET",
		Pattern:     "/account/active",
		HandlerFunc: ActivateAccountHandler,
	},
}

// POST Request
// URL : /password/forgotpassword
func ForgotPasswordHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Token Resource] Request to ForgotPasswordHandler, %s, username=%s", c.Request.URL.String(), c.PostForm("username"))
	matchToken := storage.TokenFormat{
		PasswordID: c.PostForm("username"),
	}
	updateToken := storage.TokenFormat{
		PasswordID: c.PostForm("username"),
		Token:      storage.GenerateID(time.Now().String() + "#" + c.PostForm("username")),
		CreateAt:   time.Now(),
	}
	exist, err := CheckUserExist(c.PostForm("username"))
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	if !exist {
		f.WriteHTTPError(api.ErrMailSystem, fmt.Sprint("Does not exist this mail!"))
		return
	}

	log.Infof("[Token Resource] Request Detail, match=%v, update=%v", matchToken, updateToken)
	err = storage.MongoTS.UpdateToken(matchToken, updateToken, true)
	if err != nil {
		f.WriteHTTPError(api.ErrMailSystem, fmt.Sprintf("[Token Resource] Err Msg=%v", err))
		return
	}
	log.Infof("[Token Resource] Update Token, Details=%v", updateToken)
	subject := "用户重置密码"
	curdir, _ := path.Getcurdir()
	temp := template.Must(template.ParseFiles(curdir + "/../../ds/mail/forgot.html"))
	buf := new(bytes.Buffer)
	temp.Execute(buf, struct{ Link string }{"http://" + c.Request.Host + "/resetpassword?username=" + c.PostForm("username") + "&token=" + updateToken.Token})
	mail := email.NewEmail(updateToken.PasswordID, subject, buf.String(), "text/html")
	err = mail.SendEmail(email.DeepshareDefaultEmailHost())
	if err != nil {
		f.WriteHTTPError(api.ErrMailSystem, fmt.Sprintf("[Token Resource] Send email failed! Err Msg=%v", err))
		return
	}

	c.String(http.StatusOK, "Already sent mail to user's mailbox!")
}

// POST Request
// URL : /password/resetpassword
func ResetPasswordHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	if validator.ValidateParams(c, [][]string{
		{"username", "pf"},
		{"password", "pf"},
	}) != nil {
		return
	}

	log.Infof("[Token Resource] Request to ResetPasswordHandler, %s", c.Request.URL.String())
	if c.PostForm("username") == "" || c.PostForm("password") == "" {
		f.WriteHTTPError(api.ErrBadRequestBody, "")
		return
	}

	matchToken := storage.TokenFormat{
		PasswordID: c.PostForm("username"),
	}

	result, err := storage.MongoTS.GetToken(matchToken, true)
	log.Infof("[Token Resource] Request Detail, match=%v", matchToken)
	if err != nil || len(result) == 0 || result[0].Token != c.PostForm("token") {
		f.WriteHTTPError(api.ErrTokenInvaild, "")
		return
	}
	matchUser := storage.UserFormat{
		Username: c.PostForm("username"),
	}
	updateUser := storage.UserFormat{
		Password: c.PostForm("password"),
	}
	err = storage.MongoTS.DelToken(matchToken)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	err = storage.MongoUS.UpdateUser(matchUser, updateUser, false)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{true})
}

// Get Request
// URL : /account/active
func ActivateAccountHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Token Resource] Request to ActivateAccountHandler, %s", c.Request.URL.String())
	if c.Query("username") == "" || c.Query("token") == "" {
		f.Redirect(http.StatusTemporaryRedirect, handler.LoginURLActivateFail)
		return
	}

	matchToken := storage.TokenFormat{
		AccountID: c.Query("username"),
	}

	result, err := storage.MongoTS.GetToken(matchToken, true)
	log.Infof("[Token Resource] Request Detail, match=%v result=%v err=%v", matchToken, result, err)
	if err != nil || len(result) == 0 || result[0].Token != c.Query("token") {
		log.Infof("[Token Resource] Token is nonexist or expired! Err Msg=%v", err)
		f.Redirect(http.StatusTemporaryRedirect, handler.LoginURLActivateFail)
		return
	}
	err = storage.MongoTS.DelToken(matchToken)

	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		log.Error(err)
		f.Redirect(http.StatusTemporaryRedirect, handler.LoginURLActivateFail)
		return
	}

	matchUser := storage.UserFormat{
		Username: c.Query("username"),
	}
	updateUser := storage.UserFormat{
		Activate: "1",
	}
	err = storage.MongoUS.UpdateUser(matchUser, updateUser, false)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		log.Error(err)
		f.Redirect(http.StatusTemporaryRedirect, handler.LoginURLActivateFail)
		return
	}
	session := sessions.Default(c)
	session.Set("AccountID", c.Query("username"))
	err = session.Save()
	if err != nil {
		log.Error(err)
		f.Redirect(http.StatusTemporaryRedirect, handler.LoginURLActivateFail)
		return
	}
	f.Redirect(http.StatusTemporaryRedirect, handler.LoginURLActivateSuccess)
}
