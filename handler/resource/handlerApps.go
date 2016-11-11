package resource

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/cert"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/handler/auth"
	"github.com/MISingularity/deepdash/pkg/errorutil"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/MISingularity/deepdash/pkg/sessionutil"
	"github.com/MISingularity/deepdash/pkg/util"
	"github.com/MISingularity/deepdash/pkg/validator"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
	"gopkg.in/mgo.v2/bson"
)

//var ICON_URL string = "http://7xp89t.com1.z0.glb.clouddn.com/"
var ICON_URL string = "https://nzjddxpun.qnssl.com/"

var RoutesApp = map[string]handler.Route{
	// GET Request
	"GetApplist": handler.Route{
		Name:        "GetApplist",
		Method:      "GET",
		Pattern:     "/apps",
		HandlerFunc: GetApplistHandler,
	},
	"GetAppinfo": handler.Route{
		Name:        "GetAppinfo",
		Method:      "GET",
		Pattern:     "/apps/:appid",
		HandlerFunc: GetAppinfoHandler,
	},

	// POST Request
	"PostApp": handler.Route{
		Name:        "PostApp",
		Method:      "POST",
		Pattern:     "/apps",
		HandlerFunc: PostAppHandler,
	},
	"PostCallBackUrl": handler.Route{
		Name:        "PostCallBackUrl",
		Method:      "POST",
		Pattern:     "/apps/:appid/callback",
		HandlerFunc: PostCallBackUrlHandler,
	},
	"UploadImage": handler.Route{
		Name:        "UploadImage",
		Method:      "POST",
		Pattern:     "/apps/:appid/uploadimage",
		HandlerFunc: UploadImageHandler,
	},
	"UploadIcon": handler.Route{
		Name:        "UploadIcon",
		Method:      "POST",
		Pattern:     "/apps/:appid/uploadicon",
		HandlerFunc: UploadIconHandler,
	},

	// PUT Request
	"PutApp": handler.Route{
		Name:        "ModifyApp",
		Method:      "PUT",
		Pattern:     "/apps/:appid",
		HandlerFunc: PutAppHandler,
	},

	"PutAtrributionPushUrl": handler.Route{
		Name:        "ModifyAttributionPushUrl",
		Method:      "PUT",
		Pattern:     "/apps/:appid/url",
		HandlerFunc: PutAppUrlHandler,
	},

	// DELETE Request
	"DeleteApp": handler.Route{
		Name:        "DeleteApp",
		Method:      "DELETE",
		Pattern:     "/apps/:appid",
		HandlerFunc: DeleteAppHandler,
	},
}

// GET Request
// URL : /apps
func GetApplistHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to GetApplistHandler, %s", c.Request.URL.String())
	log.Infof("[App Resource] Request Detail, AccountID : %s", sessionutil.GetAccountID(c))
	accountid := sessionutil.GetAccountID(c)
	if accountid == "" {
		f.WriteHTTPError(api.ErrBadRequestBody, "no accountid")
		return
	}

	matchBson := bson.M{"accountid": accountid}
	if auth.IsSuperAdministrator(accountid) {
		matchBson = bson.M{}
	} else if auth.IsDemoAccount(accountid) {
		// demo account can retrive app: 'linux command'
		matchBson = bson.M{"appid": "38CCA4C77072DDC9"}
	}

	result, err := storage.MongoAS.GetApplistBson(matchBson)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	applist := make([]api.AppListInfo, len(result))
	for i, _ := range result {
		applist[i].AppID = result[i].AppID
		applist[i].AppName = result[i].AppName
		applist[i].IconUrl = result[i].IconUrl
		/*
			applist[i].ChannelInfo, err = deepstats.GetAppChannels(result[i].AppID)
			if err != nil {
				f.WriteHTTPError(
					api.ErrDeepstatsConnectFailed,
					fmt.Sprintf("[App Resource] Get Deepstats app's channel list failed! Err Msg=%v", err),
				)
				return
			}*/
	}
	resdata := api.AppsGetReturnType{Applist: applist}
	log.Debugf("[App Resource][--Result][handler][GetApplistHandler] resdata=%v", resdata)
	c.JSON(http.StatusOK, resdata)
}

// GET Request
// URL : /apps/:appid
func GetAppinfoHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to GetAppinfoHandler, %s", c.Request.URL.String())
	log.Infof("[App Resource] Request Detail, AppID : %s", c.Param("appid"))
	as, err := storage.MongoAS.GetApp(storage.AppFormat{AppID: c.Param("appid")})
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, as)
}

func GenerateAppid(c *gin.Context) (string, error) {
	var appid string
	for {
		appid = storage.GenerateID(sessionutil.GetAccountID(c) + "#" + c.PostForm("name") + "#" + strconv.FormatInt(time.Now().Unix(), 10))
		// FIXME: race condition
		_, err := storage.MongoAS.GetAppBson(bson.M{"appid": appid})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				break
			}
			return "", err
		}
	}
	return appid, nil
}

// POST Request
// URL : /apps
func PostAppHandler(c *gin.Context) {
	if validator.ValidateParams(c, [][]string{
		{"appName", "pf"},
		{"pkgName", "pf"},
		{"theme", "pf"},
		{"iosBundler", "pf"},
		{"iosDownloadUrl", "pf"},
		{"iosTeamID", "pf"},
		{"iosYYBEnableBelow9", "pf"},
		{"iosYYBEnableAbove9", "pf"},
		{"androidPkgname", "pf"},
		{"androidIsDownloadDirectly", "pf"},
		{"androidHost", "pf"},
		{"androidSHA256", "pf"},
		{"androidYYBEnable", "pf"},
		{"androidDownloadUrl", "pf"},
		{"yyburl", "pf"},
		{"yybenable", "pf"},
		{"attributionPushUrl", "pf"},
		{"iconUrl", "pf"},
		{"downloadTitle", "pf"},
		{"downloadMsg", "pf"},
		{"userConfBgWeChatAndroidTipUrl", "pf"},
		{"userConfBgWeChatIosTipUrl", "pf"},
	}) != nil {
		return
	}

	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to PostAppHandler, %s", c.Request.URL.String())
	appid, err := GenerateAppid(c)
	if err != nil {
		f.WriteHTTPError(api.ErrInternalServer, "Mongo generate appid failed! Err Msg="+err.Error())
		return
	}

	seqid, err := storage.MongoSS.GetSequenceId()
	shortid := util.SequenceToShortHash(seqid, 4)

	log.Infof("[App Resource] Sequence Id is, %v, short is: %v", seqid, shortid)
	session := sessions.Default(c)
	session.Set("appid", appid)
	err = session.Save()
	if err != nil {
		f.WriteHTTPError(api.ErrInternalServer, "Save session failed! Err Msg="+err.Error())
		return
	}
	applink := "false"
	androidDownloadUrl := c.PostForm("androidDownloadUrl")
	if c.PostForm("androidSHA256") != "" {
		applink = "true"
	}
	unlink := "false"
	if c.PostForm("iosBundler") != "" && c.PostForm("iosTeamID") != "" {
		unlink = "true"
	}

	insertApp := storage.AppFormat{
		AppID:     appid,
		ShortID:   shortid,
		AppName:   c.PostForm("appName"),
		PkgName:   c.PostForm("pkgName"),
		AccountID: sessionutil.GetAccountID(c),
		Theme:     c.PostForm("theme"),

		IosBundler:             c.PostForm("iosBundler"),
		IosScheme:              "ds" + appid,
		IosDownloadUrl:         c.PostForm("iosDownloadUrl"),
		IosUniversalLinkEnable: boolPlaceholder(unlink),
		IosTeamID:              c.PostForm("iosTeamID"),
		IosYYBEnableBelow9:     boolPlaceholder(c.PostForm("iosYYBEnableBelow9")),
		IosYYBEnableAbove9:     boolPlaceholder(c.PostForm("iosYYBEnableAbove9")),

		AndroidPkgname:            c.PostForm("androidPkgname"),
		AndroidScheme:             "ds" + appid,
		AndroidIsDownloadDirectly: boolPlaceholder(c.PostForm("androidIsDownloadDirectly")),
		AndroidHost:               c.PostForm("androidHost"),
		AndroidAppLink:            boolPlaceholder(applink),
		AndroidSHA256:             c.PostForm("androidSHA256"),
		AndroidDownloadUrl:        androidDownloadUrl,
		AndroidYYBEnable:          boolPlaceholder(c.PostForm("androidYYBEnable")),

		YYBurl:             c.PostForm("yyburl"),
		YYBenable:          boolPlaceholder(c.PostForm("yybenable")),
		AttributionPushUrl: c.PostForm("attributionPushUrl"),
		IconUrl:            c.PostForm("iconUrl"),
		DownloadTitle:      c.PostForm("downloadTitle"),
		DownloadMsg:        c.PostForm("downloadMsg"),

		UserConfBgWeChatAndroidTipUrl: c.PostForm("userConfBgWeChatAndroidTipUrl"),
		UserConfBgWeChatIosTipUrl:     c.PostForm("userConfBgWeChatIosTipUrl"),
	}

	insertinfo, err := json.Marshal(insertApp)
	if err != nil {
		f.WriteHTTPError(api.ErrBadJSONBody, "")
		return
	}

	log.Infof("[App Resource] Request Detail, InsertInfo : %s", string(insertinfo))
	log.Infof("[App Resource] Ios teamid: %s", insertApp.IosTeamID)

	// append new element into apple-app-site-association-unsigned json file and sign it.
	if insertApp.IosTeamID != "" {
		log.Infof("[App Resource] Start generate new json file and sign it!")
		go cert.GenerateCertFile(insertApp.AppID, insertApp.IosTeamID, insertApp.IosBundler)
	}
	// append new element for android app link
	if insertApp.AndroidSHA256 != "" {
		log.Infof("[App Resource] Start generate new assetlinks.json file!")
		go cert.GenerateAppLinkJsonFile(insertApp.AndroidPkgname, insertApp.AndroidSHA256)
	}
	// storage phase
	// TODO:
	// Support Mongo with two phase commit,
	// synchronize mongo and redis storage server.
	err = storage.MongoAS.InsertApp(insertApp)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	err = storage.RemoteAppInfoAS.InsertApp(insertApp)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	f.WriteData(http.StatusOK, map[string]string{"appid": appid})
}

// POST Request
// URL : /apps/:appid/callback
func PostCallBackUrlHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to PostCallBackUrlHandler, %s", c.Request.URL.String())

	data := c.PostForm("data")
	url := c.PostForm("url")

	body := bytes.NewBuffer([]byte(data))
	res, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		log.Errorf("[App Resource] Failed to http post! Err Msg=%v", err)
		return
	}

	if res.StatusCode == http.StatusOK {
		f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
	} else {
		f.WriteHTTPError(api.ErrCallBackFailed, "Test failed!")
		return
	}
}

// POST Request
// URL : /apps/:appid/uploadimage
func UploadImageHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to UploadImageHandler, %s", c.Request.URL.String())
	// set the key for 七牛
	ACCESS_KEY = "KeUlmhSN6DboItUFfNss_I8Vu6M-CJP2ctmgVy1I"
	SECRET_KEY = "nxtRKNLeNjdl9tYhxAMJ6B6BBa2z1sfVPpVRPOWR"

	appid := c.Param("appid")
	file, _, e := c.Request.FormFile("uploadfile")
	if e != nil {
		log.Errorf("[App Resource] Failed to get the upload icon file! Err Msg=%v", e)
		return
	}
	defer file.Close()

	bucket := c.PostForm("bucket")
	if bucket == "" {
		bucket = "icon"
	}
	var err error
	var ret io.PutRet
	var extra = &io.PutExtra{}
	var scopekey = appid + "/" + bucket + "/" + time.Now().Format("20060102150405")
	log.Infof("[App Resource] scopekey : %s", scopekey)
	var uploadtoken = uptoken("deepshare:" + scopekey)
	log.Infof("[App Resource] Uploadtoken : %s", uploadtoken)

	// ret       	 变量用于存取返回的信息，详情见 io.PutRet
	// uploadtoken    为业务服务器端生成的上传口令
	// key       	 为文件存储的标识
	// file         	 为io.Reader类型，用于从其读取数据
	// extra     	 为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	err = io.Put(nil, &ret, uploadtoken, scopekey, file, extra)
	if err != nil {
		log.Errorf("[App Resource] Failed to upload icon file! Err Msg=%v", err)
		return
	}

	//上传成功，处理返回值
	log.Infof("[App Resource] Success to UploadIconHandler. <Hash : %s, Key : %s>", ret.Hash, ret.Key)

	// return the image url
	url := ICON_URL + scopekey
	f.WriteData(http.StatusOK, api.SimpleResponseData{Value: url})
}

// POST Request
// URL : /apps/:appid/uploadicon
func UploadIconHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to UploadIconHandler, %s", c.Request.URL.String())
	// set the key for 七牛
	ACCESS_KEY = "KeUlmhSN6DboItUFfNss_I8Vu6M-CJP2ctmgVy1I"
	SECRET_KEY = "nxtRKNLeNjdl9tYhxAMJ6B6BBa2z1sfVPpVRPOWR"

	appid := c.Param("appid")
	file, _, e := c.Request.FormFile("uploadfile")
	if e != nil {
		log.Errorf("[App Resource] Failed to get the upload icon file! Err Msg=%v", e)
		return
	}
	defer file.Close()

	var err error
	var ret io.PutRet
	var extra = &io.PutExtra{}
	var scopekey = appid + "/appicon/" + time.Now().Format("20060102150405")
	log.Infof("[App Resource] scopekey : %s", scopekey)
	var uploadtoken = uptoken("deepshare:" + scopekey)
	log.Infof("[App Resource] Uploadtoken : %s", uploadtoken)

	// ret       	 变量用于存取返回的信息，详情见 io.PutRet
	// uploadtoken    为业务服务器端生成的上传口令
	// key       	 为文件存储的标识
	// file         	 为io.Reader类型，用于从其读取数据
	// extra     	 为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	err = io.Put(nil, &ret, uploadtoken, scopekey, file, extra)
	if err != nil {
		log.Errorf("[App Resource] Failed to upload icon file! Err Msg=%v", err)
		return
	}

	//上传成功，处理返回值
	log.Infof("[App Resource] Success to UploadIconHandler. <Hash : %s, Key : %s>", ret.Hash, ret.Key)

	// return the image url
	url := ICON_URL + scopekey
	f.WriteData(http.StatusOK, api.SimpleResponseData{Value: url})

	// delete the old icon image if exist
	// get current appinfo details
	currentAppFormat, err := storage.MongoAS.GetApp(storage.AppFormat{AppID: appid})
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	log.Infof("[App Resource] currentAppFormat.Icon_Url : %s", currentAppFormat.IconUrl)
	// icon url is something like : https://nzjddxpun.qnssl.com/14293067237813bd/appicon/20151218181728
	if currentAppFormat.IconUrl != "" && !strings.Contains(currentAppFormat.IconUrl, "default") {
		oldScopeKey := strings.Replace(currentAppFormat.IconUrl, "https://nzjddxpun.qnssl.com/", "", -1)
		var rsCli = rs.New(nil)
		err = rsCli.Delete(nil, "deepshare", oldScopeKey)
		if err != nil {
			log.Errorf("[App Resource] Failed to delete old icon file! Err Msg=%v", err)
			return
		}
	}
	// 删除旧Icon成功
	log.Infof("[App Resource] Success to Delete old app icon!")
}

func uptoken(bucketName string) string {
	log.Infof("[App Resource] uptoken bucketName : %s", bucketName)
	putPolicy := rs.PutPolicy{
		Scope: bucketName,
	}
	return putPolicy.Token(nil)
}

// PUT Reqeust
// URL : /apps/:appid/url
func PutAppUrlHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to PutAppUrlHandler, %s", c.Request.URL.String())

	appId := c.Param("appid")
	attributionPushUrl := c.PostForm("attributionpushurl")

	appFormat, err := storage.MongoAS.GetApp(storage.AppFormat{AppID: appId})
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	// insert new appformat
	updateApp := storage.AppFormat(appFormat)
	updateApp.AttributionPushUrl = attributionPushUrl

	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	// storage phase
	// TODO:
	// Support Mongo with two phase commit,
	// synchronize mongo and redis storage server.
	err = storage.RemoteAppInfoAS.UpdateApp(storage.AppFormat{}, updateApp, true)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	err = storage.MongoAS.UpdateAppBson(bson.M{"appid": appId}, bson.M{"attributionpushurl": attributionPushUrl}, true)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

// PUT Reqeust
// URL : /apps/:appid
func PutAppHandler(c *gin.Context) {
	if validator.ValidateParams(c, [][]string{
		{"appid", "p"},
		{"appName", "pf"},
		{"pkgName", "pf"},
		{"theme", "pf"},
		{"iosBundler", "pf"},
		{"iosDownloadUrl", "pf"},
		{"iosTeamID", "pf"},
		{"iosYYBEnableBelow9", "pf"},
		{"iosYYBEnableAbove9", "pf"},
		{"androidPkgname", "pf"},
		{"androidIsDownloadDirectly", "pf"},
		{"androidHost", "pf"},
		{"androidSHA256", "pf"},
		{"androidDownloadUrl", "pf"},
		{"androidYYBEnable", "pf"},
		{"yyburl", "pf"},
		{"yybenable", "pf"},
		{"attributionPushUrl", "pf"},
		{"iconUrl", "pf"},
		{"downloadTitle", "pf"},
		{"downloadMsg", "pf"},
		{"userConfBgWeChatAndroidTipUrl", "pf"},
		{"userConfBgWeChatIosTipUrl", "pf"},
	}) != nil {
		return
	}

	f := httputil.NewGinframework(c)

	log.Infof("[App Resource] Request to PutAppHandler, %s", c.Request.URL.String())
	applink := "false"
	if c.PostForm("androidSHA256") != "" {
		applink = "true"
	}
	unlink := "false"
	if c.PostForm("iosBundler") != "" && c.PostForm("iosTeamID") != "" {
		unlink = "true"
	}

	updateAppBson := validator.ParseParams(c, [][]string{
		{"appid", "p", "appid"},
		{"appName", "pf", "appname"},
		{"pkgName", "pf", "fullpkgname"},
		{"theme", "pf", "theme"},

		{"iosBundler", "pf", "iosbundler"},
		{"iosDownloadUrl", "pf", "iosdownloadurl"},
		{"iosTeamID", "pf", "iosteamid"},
		{"iosYYBEnableBelow9", "pf", "iosyybenablebelow9"},
		{"iosYYBEnableAbove9", "pf", "iosyybenableabove9"},
		{"forceDownload", "pf", "forcedownload"},

		{"androidPkgname", "pf", "androidpkgname"},
		{"androidIsDownloadDirectly", "pf", "androidisdownloaddirectly"},
		{"androidHost", "pf", "androidhost"},
		{"androidSHA256", "pf", "androidsha256"},
		{"androidDownloadUrl", "pf", "androiddownloadurl"},
		{"androidYYBEnable", "pf", "androidyybenable"},

		{"yyburl", "pf", "yyburl"},
		{"yybenable", "pf", "yybenable"},
		{"attributionPushUrl", "pf", "attributionpushurl"},
		{"iconUrl", "pf", "iconurl"},
		{"downloadTitle", "pf", "download_title"},
		{"downloadMsg", "pf", "download_msg"},

		{"userConfBgWeChatAndroidTipUrl", "pf", "userconfbgwechatandroidtipurl"},
		{"userConfBgWeChatIosTipUrl", "pf", "userconfbgwechatiostipurl"},
	})
	updateAppBson["yybenable"] = boolPlaceholder(updateAppBson["yybenable"].(string))
	updateAppBson["androidyybenable"] = boolPlaceholder(updateAppBson["androidyybenable"].(string))
	updateAppBson["androidapplink"] = boolPlaceholder(applink)
	updateAppBson["androidisdownloaddirectly"] = boolPlaceholder(updateAppBson["androidisdownloaddirectly"].(string))
	updateAppBson["iosyybenableabove9"] = boolPlaceholder(updateAppBson["iosyybenableabove9"].(string))
	updateAppBson["iosyybenablebelow9"] = boolPlaceholder(updateAppBson["iosyybenablebelow9"].(string))
	updateAppBson["forcedownload"] = boolPlaceholder(updateAppBson["forcedownload"].(string))
	updateAppBson["iosunilink"] = boolPlaceholder(unlink)

	updateAppBson["iosscheme"] = "ds" + updateAppBson["appid"].(string)
	updateAppBson["androidscheme"] = "ds" + updateAppBson["appid"].(string)

	// get current appinfo details
	currentAppFormat, err := storage.MongoAS.GetApp(storage.AppFormat{AppID: c.Param("appid")})
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	if currentAppFormat.ShortID == "" {
		seqid, _ := storage.MongoSS.GetSequenceId()
		shortid := util.SequenceToShortHash(seqid, 4)
		updateAppBson["shortid"] = shortid
	}

	upsertinfo, err := json.Marshal(updateAppBson)
	if err != nil {
		f.WriteHTTPError(api.ErrBadJSONBody, "")
		return
	}
	log.Infof("[App Resource] Request Detail, UpsertInfo : %s", string(upsertinfo))

	err = storage.MongoAS.UpdateAppBson(bson.M{"appid": c.Param("appid")}, updateAppBson, false)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})

	// append new element into apple-app-site-association-unsigned json file and sign it.
	if updateAppBson["iosbundler"] != currentAppFormat.IosBundler && updateAppBson["iosbundler"] != "" ||
		updateAppBson["iosteamid"] != currentAppFormat.IosTeamID && updateAppBson["iosteamid"] != "" ||
		currentAppFormat.ShortID == "" {
		log.Infof("[App Resource] Start generate new json file and sign it!")
		go cert.GenerateCertFile(updateAppBson["appid"].(string), updateAppBson["iosteamid"].(string), updateAppBson["iosbundler"].(string))
	}
	// append new element for android app link
	if updateAppBson["androidsha256"] != currentAppFormat.AndroidSHA256 && updateAppBson["androidsha256"] != "" {
		log.Infof("[App Resource] Start generate new assetlinks.json file!")
		go cert.GenerateAppLinkJsonFile(updateAppBson["androidpkgname"].(string), updateAppBson["androidsha256"].(string))
	}

	// storage phase
	// mannually update and push to redis
	_ = json.Unmarshal([]byte(upsertinfo), &currentAppFormat)
	log.Debugf("[App Resource] update redis: %v", currentAppFormat)

	err = storage.RemoteAppInfoAS.UpdateApp(storage.AppFormat{}, currentAppFormat, true)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
}

func DeleteAppHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	if validator.ValidateParams(c, [][]string{
		{"appid", "p"},
	}) != nil {
		return
	}

	err := storage.MongoAS.DeleteApp(c.Param("appid"))
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

func boolPlaceholder(boolv string) string {
	if boolv != "true" {
		return "false"
	}
	return "true"
}
