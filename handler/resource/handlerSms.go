package resource

import (
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/pkg/errorutil"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/validator"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/gin"
)

func GetSmsListHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	appID := c.Param("appid")
	match := bson.M{"appid": appID}

	smses, _ := storage.MongoSmsService.GetSmsList(match)
	f.WriteData(http.StatusOK, api.Smses{Smses: smses})
}

func PostSmsHandler(c *gin.Context) {
	if validator.ValidateParams(c, [][]string{
		{"appid", "p"},
		{"url", "pf"},
		{"content", "pf"},
	}) != nil {
		return
	}

	f := httputil.NewGinframework(c)
	appId := c.Param("appid")
	channelName := c.PostForm("channelname")
	url := c.PostForm("url")
	content := c.PostForm("content")

	err := storage.MongoSmsService.InsertSms(bson.M{
		"appid":       appId,
		"channelname": channelName,
		"url":         url,
		"content":     content,
		"createat":    strconv.FormatInt(time.Now().Unix(), 10),
	})

	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

func DeleteSmsHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	appId := c.Param("appid")
	smsId := c.Param("smsid")

	err := storage.MongoSmsService.RemoveSms(bson.M{"_id": bson.ObjectIdHex(smsId), "appid": appId})
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

func UpdateSmsHandler(c *gin.Context) {
	if validator.ValidateParams(c, [][]string{
		{"appid", "p"},
		{"content", "pf"},
	}) != nil {
		return
	}

	f := httputil.NewGinframework(c)
	appId := c.Param("appid")
	smsId := c.Param("smsid")
	content := c.PostForm("content")

	err := storage.MongoSmsService.UpdateSms(appId, smsId, content)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}
