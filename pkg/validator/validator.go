package validator

import (
	"net/http"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

var ParamConf = map[string]string{
	"appid":                         "min=1,max=50",
	"appName":                       "min=1,max=50",
	"pkgName":                       "max=300",
	"theme":                         "max=3,regexp=^[0-9]*$",
	"iosBundler":                    "max=100",
	"iosDownloadUrl":                "max=1000",
	"iosTeamID":                     "max=20",
	"iosYYBEnableBelow9":            "max=300",
	"iosYYBEnableAbove9":            "max=300",
	"androidPkgname":                "max=300",
	"androidIsDownloadDirectly":     "max=300",
	"androidHost":                   "max=300",
	"androidSHA256":                 "max=300",
	"androidDownloadUrl":            "max=1000",
	"androidYYBEnable":              "max=300",
	"yyburl":                        "max=1000",
	"yybenable":                     "max=300",
	"attributionPushUrl":            "max=1000",
	"iconUrl":                       "max=1000",
	"downloadTitle":                 "max=500",
	"downloadMsg":                   "max=1000",
	"userConfBgWeChatAndroidTipUrl": "max=1000",
	"userConfBgWeChatIosTipUrl":     "max=1000",

	"username": "max=50",
	"password": "max=50",
	"phone":    "max=15,regexp=^[0-9]*$",
	"source":   "max=200",
	"appname":  "max=50",

	// sms
	"url":     "max=50",
	"content": "max=100",
}

//"phone":    "regexp=^()|([0-9]{7, 15})$",
func ValidateParams(c *gin.Context, keys [][]string) error {
	f := httputil.NewGinframework(c)
	for _, param := range keys {
		var value = ""
		key := param[0]
		src := param[1]
		if checker, ok := ParamConf[key]; ok {
			if src == "q" {
				value = c.Query(key)
			} else if src == "p" {
				value = c.Param(key)
			} else if src == "pf" {
				value = c.PostForm(key)
			} else {
				value = c.Query(key)
			}
			log.Debugf("name: %s, value: %s, checker: %s", key, value, checker)
			err := validator.Valid(value, checker)
			if err != nil {
				log.Errorf("Validate params error, %s: %s", key, err)

				f.WriteHTTPError(
					httputil.HTTPError{
						StatusCode: http.StatusBadRequest,
						Code:       api.CodeParamNoneValid,
						Message:    key + ": " + err.Error(),
						Fatal:      false},
					err.Error())
				return err
			}
		} else {
			log.Errorf("Validate key not exist error, %s", key)
		}
	}
	return nil
}

func ParseParams(c *gin.Context, keys [][]string) bson.M {
	// TODO: key lost in query
	res := bson.M{}
	for _, param := range keys {
		var value = ""
		key := param[0]
		src := param[1]
		dst := param[2]

		if src == "q" {
			value = c.Query(key)
		} else if src == "p" {
			value = c.Param(key)
		} else if src == "pf" {
			value = c.PostForm(key)
		} else {
			value = c.Query(key)
		}
		log.Debugf("name: %s, value: %s", key, value)
		res[dst] = value
	}
	return res
}
