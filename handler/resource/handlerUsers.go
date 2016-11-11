package resource

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/handler/client"
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

var RoutesRegister = map[string]handler.Route{
	"TemporaryPutUser": handler.Route{
		Name:        "PutUser",
		Method:      "PUT",
		Pattern:     "/202cb962ac59075b964b07152d234b70/users",
		HandlerFunc: PutUserHandler,
	},
	"PostUser": handler.Route{
		Name:        "PostUser",
		Method:      "POST",
		Pattern:     "/this-is-a-clandestine-resource/users",
		HandlerFunc: PostUserHandler,
	},
	"ApplyTry": handler.Route{
		Name:        "ApplyTry",
		Method:      "POST",
		Pattern:     "/apply-try/users",
		HandlerFunc: ApplyTryHandler,
	},
	"ResentEmail": handler.Route{
		Name:        "ResentEmail",
		Method:      "POST",
		Pattern:     "/this-is-a-clandestine-resource/resent-email",
		HandlerFunc: ResentEmail,
	},
}

func ResentEmail(c *gin.Context) {

	f := httputil.NewGinframework(c)
	log.Infof("[User Resource] Request to ResentEmail, %s", c.Request.URL.String())

	session := sessions.Default(c)
	emailObj := session.Get("unverified-email")

	if emailObj == nil {
		f.WriteHTTPError(api.ErrBadRequestBody, "[User Resource] Lack required params")
		return
	}

	// check username format
	if validator.ValidateParams(c, [][]string{
		{"username", "pf"},
	}) != nil {
		return
	}

	if c.PostForm("username") == "" {
		f.WriteHTTPError(api.ErrBadRequestBody, "[User Resource] Lack required params")
		return
	}

	// check if username is not active
	notActive, user, _, err := client.CheckUserIsNotActive(c.PostForm("username"), "")
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	if notActive == false {
		if err != nil {
			f.WriteData(http.StatusOK, api.ResponseData{Success: false, Message: "服务器连接失败，请重试!"})
		} else if user.Username == "" {
			f.WriteData(http.StatusOK, api.ResponseData{Success: false, Message: "用户不存在!"})
		} else {
			f.WriteData(http.StatusOK, api.ResponseData{Success: false, Message: "用户已激活!"})
		}
		return
	}

	// update old token
	matchToken := storage.TokenFormat{
		AccountID: c.PostForm("username"),
	}

	updateToken := storage.TokenFormat{
		AccountID: c.PostForm("username"),
		Token:     storage.GenerateID(time.Now().String() + "#" + c.PostForm("username")),
		CreateAt:  time.Now(),
	}

	err = storage.MongoTS.UpdateToken(matchToken, updateToken, true)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	log.Info("[Token Resource] Update Token, Details=%v", updateToken)

	// write user mail content
	subject := "欢迎使用deepshare"

	curdir, _ := path.Getcurdir()
	temp := template.Must(template.ParseFiles(curdir + "/../../ds/mail/register.html"))

	buf := new(bytes.Buffer)
	temp.Execute(buf, struct{ Link string }{"http://" + c.Request.Host + "/account/active?username=" + c.PostForm("username") + "&token=" + updateToken.Token})

	// send mail to user
	mail := email.NewEmail(c.PostForm("username"), subject, buf.String(), "text/html")
	err = mail.SendEmail(email.DeepshareDefaultEmailHost())
	if err != nil {
		f.WriteHTTPError(api.ErrMailSystem, fmt.Sprintf("[User Resource] Send email failed! Err Msg=%v", err))
		return
	}

	// write bd mail content
	message := "App 名称 : " + user.Appname + ";\n"
	message += "邮箱地址 : " + user.Username + ";\n"
	message += "联系电话 : " + user.Phone + ";\n"
	message += "来源 : " + user.Source + ";\n"

	log.Infof("[User Resource] Email message --> %s", message)

	// send email to bd
	go email.NewEmail(email.RECEIVER, "用户重发激活邮件", message, "text/plain").SendEmail(email.DeepshareDefaultEmailHost())

	log.Info("[User Resource] resend email success!")
	f.WriteData(http.StatusOK, api.ResponseData{Success: true, Message: "邮件发送成功，请在15分钟内激活"})
}

// POST Reqeust
// URL : /this-is-a-clandestine-resource/users
func PostUserHandler(c *gin.Context) {
	if validator.ValidateParams(c, [][]string{
		{"username", "pf"},
		{"password", "pf"},
		{"appname", "pf"},
		{"phone", "pf"},
		{"source", "pf"},
	}) != nil {
		return
	}

	f := httputil.NewGinframework(c)
	log.Infof("[User Resource] Request to PostUserHandler, %s", c.Request.URL.String())

	if c.PostForm("username") == "" || c.PostForm("password") == "" {
		f.WriteHTTPError(api.ErrBadRequestBody, "[User Resource] Lack required params")
		return
	}

	exist, err := client.CheckUserExist(c.PostForm("username"))
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	if exist {
		f.WriteHTTPError(api.ErrUserAlreadyExist, "Account already exists!")
		return
	}

	if !handler.DEBUG {
		matchToken := storage.TokenFormat{
			AccountID: c.PostForm("username"),
		}
		updateToken := storage.TokenFormat{
			AccountID: c.PostForm("username"),
			Token:     storage.GenerateID(time.Now().String() + "#" + c.PostForm("username")),
			CreateAt:  time.Now(),
		}
		err = storage.MongoTS.UpdateToken(matchToken, updateToken, true)
		httperror = errorutil.ProcessMongoError(f, err)
		if httperror {
			return
		}
		log.Info("[Token Resource] Update Token, Details=%v", updateToken)

		subject := "欢迎使用deepshare"

		curdir, _ := path.Getcurdir()
		temp := template.Must(template.ParseFiles(curdir + "/../../ds/mail/register.html"))

		buf := new(bytes.Buffer)
		temp.Execute(buf, struct{ Link string }{"http://" + c.Request.Host + "/account/active?username=" + c.PostForm("username") + "&token=" + updateToken.Token})

		mail := email.NewEmail(c.PostForm("username"), subject, buf.String(), "text/html")
		err = mail.SendEmail(email.DeepshareDefaultEmailHost())
		if err != nil {
			f.WriteHTTPError(api.ErrMailSystem, fmt.Sprintf("[User Resource] Send email failed! Err Msg=%v", err))
			return
		}
	}

	title := "新用户注册 - " + c.PostForm("appname")
	message := "App 名称 : " + c.PostForm("appname") + ";\n"
	message += "邮箱地址 : " + c.PostForm("username") + ";\n"
	message += "联系电话 : " + c.PostForm("phone") + ";\n"
	message += "来源 : " + c.PostForm("source") + ";\n"

	log.Infof("[User Resource] Email message --> %s", message)

	go email.NewEmail(email.RECEIVER, title, message, "text/plain").SendEmail(email.DeepshareDefaultEmailHost())

	matchUser := storage.UserFormat{
		Username: c.PostForm("username"),
	}
	insertUser := storage.UserFormat{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Phone:    c.PostForm("phone"),
		Source:   c.PostForm("source"),
		Appname:  c.PostForm("appname"),
		CreateAt: strconv.FormatInt(time.Now().Unix(), 10),
		Activate: "0",
	}
	err = storage.MongoUS.UpdateUser(matchUser, insertUser, true)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	log.Info("[User Resource] Register success!")
	f.WriteData(http.StatusOK, api.SimpleResponseData{Value: "Already registered! Please activate your account in 15 minutes, otherwise your account will be removed!"})

}

// POST Request
func ApplyTryHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[User Resource] Request to ApplyTryHandler, %s", c.Request.URL.String())
	message := "App 名称 : " + c.PostForm("appname") + ";\n"
	message += "App 日活量 : " + c.PostForm("appactive") + ";\n"
	message += "iOS 下载地址 : " + c.PostForm("iosdownloadurl") + ";\n"
	message += "Android 下载地址 : " + c.PostForm("androiddownloadurl") + ";\n"
	message += "邮箱地址 : " + c.PostForm("emailaddress") + ";\n"
	message += "联系电话 : " + c.PostForm("phonenumber") + ";\n"
	message += "来源 : " + c.PostForm("source") + ";\n"

	log.Infof("[User Resource] Email message --> %s", message)

	err := email.NewEmail(email.RECEIVER, "新用户申请试用", message, "text/plain").SendEmail(email.DeepshareDefaultEmailHost())
	if err != nil {
		f.WriteHTTPError(api.ErrMailSystem, "")
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

// PUT Request
// URL : /this-is-a-clandestine-resource/users
func PutUserHandler(c *gin.Context) {
	if validator.ValidateParams(c, [][]string{
		{"username", "pf"},
		{"password", "pf"},
		{"appname", "pf"},
		{"phone", "pf"},
		{"source", "pf"},
	}) != nil {
		return
	}

	f := httputil.NewGinframework(c)
	log.Infof("[User Resource] Request to PutUserHandler, %s", c.Request.URL.String())

	if c.PostForm("username") == "" || c.PostForm("password") == "" || !strings.Contains(c.PostForm("username"), "@") {
		f.WriteHTTPError(api.ErrBadRequestBody, "[User Resource] Lack required params or incorrect param")
		return
	}
	exist, err := client.CheckUserExist(c.PostForm("username"))
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	if exist {
		f.WriteHTTPError(api.ErrUserAlreadyExist, "Account already exists!")
		return
	}
	match := storage.UserFormat{
		Username: c.PostForm("username"),
	}
	updateUser := storage.UserFormat{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Phone:    c.PostForm("phone"),
		Source:   c.PostForm("source"),
		Appname:  c.PostForm("appname"),
		CreateAt: strconv.FormatInt(time.Now().Unix(), 10),
		Activate: "1",
	}
	log.Infof("[User Resource] Request Detail, match=%v, update=%v", match, updateUser, true)
	err = storage.MongoUS.UpdateUser(match, updateUser, true)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}
