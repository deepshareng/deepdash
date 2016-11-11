package private

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/deepstats"
	"github.com/MISingularity/deepdash/pkg/errorutil"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/privateutil"
	"github.com/MISingularity/deepdash/pkg/transmit"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

var RoutesPrivateResource = map[string]handler.Route{
	"GetStatus": handler.Route{
		Name:        "GetUserStatus",
		Method:      "GET",
		Pattern:     "/this-is-a-clandestine-resource/status",
		HandlerFunc: GetStatusHandler,
	},
	"GetTotalStatus": handler.Route{
		Name:        "GetTotalStatus",
		Method:      "GET",
		Pattern:     "/this-is-a-clandestine-resource/total-status",
		HandlerFunc: GetTotalStatusHandler,
	},
	"GetSpecificAppStatus": handler.Route{
		Name:        "GetSpecificAppStatus",
		Method:      "GET",
		Pattern:     "/this-is-a-clandestine-resource/apps",
		HandlerFunc: GetSpecificAppStatus,
	},

	"GetUserList": handler.Route{
		Name:        "GetUserList",
		Method:      "GET",
		Pattern:     "/this-is-a-clandestine-resource/users",
		HandlerFunc: GetUserList,
	},

	// POST Request
	"FreezeUserById": handler.Route{
		Name:        "FreezeUserById",
		Method:      "POST",
		Pattern:     "/this-is-a-clandestine-resource/freeze",
		HandlerFunc: FreezeUser,
	},
	"UnfreezeUserById": handler.Route{
		Name:        "UnfreezeUserById",
		Method:      "POST",
		Pattern:     "/this-is-a-clandestine-resource/unfreeze",
		HandlerFunc: UnfreezeUser,
	},
	"UpdateReturnVisitById": handler.Route{
		Name:        "UpdateReturnVisitById",
		Method:      "POST",
		Pattern:     "/this-is-a-clandestine-resource/update/return-visit",
		HandlerFunc: UpdateReturnVisitById,
	},
	"UserAdd": handler.Route{
		Name:        "UserAdd",
		Method:      "POST",
		Pattern:     "/this-is-a-clandestine-resource/user/add",
		HandlerFunc: UserAdd,
	},
	"PermissionAdd": handler.Route{
		Name:        "PermissionAdd",
		Method:      "POST",
		Pattern:     "/this-is-a-clandestine-resource/permission/add",
		HandlerFunc: PermissionAdd,
	},
	"PermissionRemove": handler.Route{
		Name:        "PermissionRemove",
		Method:      "POST",
		Pattern:     "/this-is-a-clandestine-resource/permission/remove",
		HandlerFunc: PermissionRemove,
	},
}

func PermissionRemove(c *gin.Context) {
	f := httputil.NewGinframework(c)
	userId := c.PostForm("id")
	permission := c.PostForm("permission")
	userIdBson := bson.ObjectIdHex(userId)
	userResults, err := storage.MongoUS.GetUser(storage.UserFormat{Id: userIdBson}, false, []string{}, []string{})
	if err != nil || len(userResults) < 1 {
		f.WriteHTTPError(api.ErrInternalServer, err.Error())
		return
	}
	userPermissionList := strings.Split(userResults[0].Permission, "_")
	newPermission := ""
	for i := 0; i < len(userPermissionList); i++ {
		if userPermissionList[i] == permission {
			continue
		}
		if newPermission != "" {
			newPermission += "_"
		}
		newPermission += userPermissionList[i]
	}
	if newPermission == "" {
		newPermission = "normalUser"
	}
	matchUser := storage.UserFormat{
		Id: userIdBson,
	}
	insertUser := storage.UserFormat{
		Permission: newPermission,
	}
	err = storage.MongoUS.UpdateUser(matchUser, insertUser, false)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.ResponseData{Success: true, Message: "移除权限成功!"})
}

func PermissionAdd(c *gin.Context) {
	f := httputil.NewGinframework(c)
	userId := c.PostForm("id")
	permission := c.PostForm("permission")
	userIdBson := bson.ObjectIdHex(userId)
	userResults, err := storage.MongoUS.GetUser(storage.UserFormat{Id: userIdBson}, false, []string{}, []string{})
	if err != nil || len(userResults) < 1 {
		f.WriteHTTPError(api.ErrInternalServer, err.Error())
		return
	}
	userPermissionList := strings.Split(userResults[0].Permission, "_")
	newPermission := ""
	for i := 0; i < len(userPermissionList); i++ {
		if userPermissionList[i] == permission {
			continue
		}
		if newPermission != "" {
			newPermission += "_"
		}
		newPermission += userPermissionList[i]
	}
	if newPermission != "" {
		newPermission += "_"
	}
	newPermission += permission
	matchUser := storage.UserFormat{
		Id: userIdBson,
	}
	insertUser := storage.UserFormat{
		Permission: newPermission,
	}
	err = storage.MongoUS.UpdateUser(matchUser, insertUser, false)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.ResponseData{Success: true, Message: "添加权限成功!"})
}

func GetUserList(c *gin.Context) {
	f := httputil.NewGinframework(c)

	userResults, err := storage.MongoUS.GetUser(storage.UserFormat{Activate: "1"}, true, []string{}, []string{})
	if err != nil {
		f.WriteHTTPError(api.ErrInternalServer, err.Error())
		return
	}
	f.WriteData(http.StatusOK, userResults)
}

func UserAdd(c *gin.Context) {
	f := httputil.NewGinframework(c)
	user := storage.UserFormat{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Appname:  c.PostForm("appname"),
		Phone:    c.PostForm("phone"),
		Activate: "1",
		CreateAt: strconv.FormatInt(time.Now().Unix(), 10),
	}

	if user.Username == "" || user.Password == "" || user.Appname == "" || user.Phone == "" {
		f.WriteData(http.StatusOK, api.ResponseData{Success: false, Message: "添加用户失败，信息不完整!"})
		return
	}
	err := storage.MongoUS.InsertUser(user)

	if err != nil {
		log.Error("Insert user {username :", user.Username, ", password: ", user.Password, ", appname: ", user.Appname, ", phone: ", user.Phone, "} failed! error :", err)
		f.WriteData(http.StatusOK, api.ResponseData{Success: false, Message: "添加用户失败!"})
		return
	}
	f.WriteData(http.StatusOK, api.ResponseData{Success: true, Message: "添加用户成功!"})
}

func UpdateReturnVisitById(c *gin.Context) {

	f := httputil.NewGinframework(c)

	userId := c.PostForm("id")
	returnVisitResult := c.PostForm("returnVisitResult")
	returnVisitTime := c.PostForm("returnVisitTime")

	if err := validator.Valid(returnVisitResult, "nonzero"); err != nil {
		returnVisitResult = "Empty ^_^"
	}

	if err := validator.Valid(returnVisitTime, "nonzero"); err != nil {
		f.WriteHTTPError(api.ErrBadJSONBody, "[User Resource] Lack of the return visit data")
		return
	}

	log.Infof("[User Resource] Request to update user return visit, %s", c.Request.URL.String())

	if err := validator.Valid(userId, "nonzero"); err != nil {
		f.WriteHTTPError(api.ErrBadJSONBody, "[User Resource] Lack User Id")
		return
	}
	userIdBson := bson.ObjectIdHex(userId)
	matchUser := storage.UserFormat{
		Id: userIdBson,
	}
	insertUser := storage.UserFormat{
		ReturnVisitResult: returnVisitResult,
		ReturnVisitTime:   returnVisitTime,
	}

	err := storage.MongoUS.UpdateUser(matchUser, insertUser, false)
	httperror := errorutil.ProcessMongoError(f, err)

	if httperror {
		return
	}

	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

func FreezeUser(c *gin.Context) {
	f := httputil.NewGinframework(c)
	userId := c.PostForm("id")
	log.Infof("[User Resource] Request to FreezeUser, %s", c.Request.URL.String())
	if userId == "" {
		f.WriteHTTPError(api.ErrBadRequestBody, "[User Resource] Lack User Id")
		return
	}
	userIdBson := bson.ObjectIdHex(userId)
	matchUser := storage.UserFormat{
		Id: userIdBson,
	}
	insertUser := storage.UserFormat{
		Freeze: "1",
	}

	err := storage.MongoUS.UpdateUser(matchUser, insertUser, false)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

func UnfreezeUser(c *gin.Context) {
	f := httputil.NewGinframework(c)
	userId := c.PostForm("id")
	log.Infof("[User Resource] Request to UnFreezeUser, %s", c.Request.URL.String())
	if userId == "" {
		f.WriteHTTPError(api.ErrBadRequestBody, "[User Resource] Lack User Id")
		return
	}
	matchUser := storage.UserFormat{
		Id: bson.ObjectIdHex(userId),
	}
	insertUser := storage.UserFormat{
		Freeze: "0",
	}
	err := storage.MongoUS.UpdateUser(matchUser, insertUser, false)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

func GetSpecificAppStatus(c *gin.Context) {
	f := httputil.NewGinframework(c)
	session := sessions.Default(c)
	session.Set("appid", c.Query("appid"))
	err := session.Save()
	if err != nil {
		log.Infof("[Client] Save Session Failed, err=%v", err)
		f.WriteHTTPError(api.ErrBadRequestBody, "")
		return
	}
	f.Redirect(http.StatusTemporaryRedirect, "/")
}

func StatusCheck(event []map[string]string, allEvent []map[string]string, expiry bool) string {
	finish := true
	nonexist := true
	urlValue, err := strconv.Atoi(allEvent[0]["url-tem-overall-value"])
	if err != nil {
		log.Infof("[App ERROR] Value is Not String, ")
	}
	shareLink, err := strconv.Atoi(allEvent[0]["sharelink-tem-overall-value"])
	if err != nil {
		log.Infof("[App ERROR] Value is Not String, ")
	}

	if urlValue < 100 && shareLink < 100 {
		finish = false
	}

	for i := 0; i < 7; i++ {
		if event[i]["url-tem-overall-value"] != "0" || event[i]["sharelink-tem-overall-value"] != "0" {
			nonexist = false
		}
	}

	if finish && (!nonexist) {
		return "集成完毕"
	}

	if nonexist && expiry {
		return "集成失败"
	}

	return "集成中"
}

func UpdateAccountStatus(origin string, merge string) string {
	if origin == "集成完毕" || merge == "集成完毕" {
		return "集成完毕"
	}
	if origin == "集成中" || merge == "集成中" {
		return "集成中"
	}
	if origin == "集成失败" || merge == "集成失败" {
		return "集成失败"
	}
	return "注册完毕"
}

func EventAddition(currentCount int, addition string) (int, error) {
	if addition == "" {
		return currentCount, nil
	}
	increment, err := strconv.Atoi(addition)
	if err != nil {
		return 0, err
	}
	return currentCount + increment, nil
}

func GetStatusHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to GetStatusHandler, %s", c.Request.URL.String())

	// get user status
	userMapping := make(map[string]*api.PrivateUserStatus)
	userResults, err := storage.MongoUS.GetUser(storage.UserFormat{Activate: "1"}, true, []string{}, []string{})
	if err != nil {
		f.WriteHTTPError(api.ErrInternalServer, err.Error())
		return
	}
	for _, v := range userResults {
		userMapping[v.Id.Hex()] = &api.PrivateUserStatus{UserFormat: v, AccountStatus: "注册完毕", LastWeekAccountStatus: "注册完毕"}
	}

	// get app status
	appResults, err := storage.MongoAS.GetApplist(storage.AppFormat{AccountID: ""})
	if err != nil {
		f.WriteHTTPError(api.ErrInternalServer, err.Error())
		return
	}

	// get app data
	// user/app status update
	appStatus := []api.PrivateAppStatus{}
	userStatus := []api.PrivateUserStatus{}
	lastWeekCount := api.StatusCounter{}
	// accountCount := api.StatusCounter{}
	appAccountCount := api.StatusCounter{}
	{
		for _, v := range appResults {
			if userMapping[v.AccountID] == nil || userMapping[v.AccountID].Username == "lihai" {
				continue
			}
			hasIOS := v.IosBundler != "" || v.IosDownloadUrl != "" || v.IosTeamID != ""
			hasAndroid := v.AndroidPkgname != ""

			rawEvents := []string{"basic-event", "url-tem-overall-value", "sharelink-tem-overall-value"}
			convertedEvents := transmit.RequisiteAttrs(rawEvents, "d")
			aggregateResult, err := deepstats.GetChannelCounters(v.AppID, "all", convertedEvents, "d", strconv.FormatInt(time.Now().Add(-time.Hour*24*14).Unix(), 10), strconv.FormatInt(time.Now().Unix(), 10), "", "")

			if err != nil {
				f.WriteHTTPError(api.ErrInternalServer, err.Error())
				return
			}

			allRawEvents := []string{"basic-event", "url-tem-overall-value", "sharelink-tem-overall-value"}
			allConvertedEvents := transmit.RequisiteAttrs(allRawEvents, "t")
			allAggregateResult, err := deepstats.GetChannelCounters(v.AppID, "all", allConvertedEvents, "t", "", "", "1", "")
			allEventRes := transmit.Calculate(allRawEvents, allConvertedEvents, allAggregateResult, "t", 1)

			if err != nil {
				f.WriteHTTPError(api.ErrInternalServer, err.Error())
				return
			}

			eventRes := transmit.Calculate(rawEvents, convertedEvents, aggregateResult, "d", 14)
			currentApp := api.PrivateAppStatus{
				UserName:         userMapping[v.AccountID].Username,
				SubjectApp:       userMapping[v.AccountID].Appname,
				AppName:          v.AppName,
				AppID:            v.AppID,
				LinkDemontration: 0,
				LinkShare:        0,
				AppInstall:       0,
				AppOpen:          0,
				AccountStatus:    "注册完毕",
			}
			createat, err := strconv.ParseInt(userMapping[v.AccountID].CreateAt, 10, 64)

			if err != nil {
				createat = 0
			}

			if hasIOS || hasAndroid {
				acountStatus := StatusCheck(eventRes[:7], allEventRes, time.Now().Add(-time.Hour*24*7*2).Unix() > createat)
				userMapping[v.AccountID].AccountStatus = UpdateAccountStatus(userMapping[v.AccountID].AccountStatus, acountStatus)
				userMapping[v.AccountID].LastWeekAccountStatus = UpdateAccountStatus(userMapping[v.AccountID].LastWeekAccountStatus, StatusCheck(eventRes[7:], allEventRes, time.Now().Add(-time.Hour*24*7*3).Unix() > createat))
				currentApp.AccountStatus = UpdateAccountStatus(currentApp.AccountStatus, acountStatus)
			}

			if userMapping[v.AccountID].Freeze == "1" {
				currentApp.AppName += "(冻结)"
				currentApp.AccountStatus = "冻结中"
			}

			eventRes = eventRes[:7]
			for i := 0; i < 7; i++ {
				currentApp.LinkDemontration, _ = EventAddition(currentApp.LinkDemontration, eventRes[i]["url-tem-overall-value"])
				currentApp.LinkShare, _ = EventAddition(currentApp.LinkShare, eventRes[i]["sharelink-tem-overall-value"])
				currentApp.AppInstall, _ = EventAddition(currentApp.AppInstall, eventRes[i]["match/install_with_params"])
				currentApp.AppOpen, _ = EventAddition(currentApp.AppOpen, eventRes[i]["match/open_with_params"])
			}
			appStatus = append(appStatus, currentApp)
		}
		// conclude user status
		for _, v := range userMapping {
			if v.Freeze == "1" {
				v.AccountStatus = "冻结中"
				v.LastWeekAccountStatus = "冻结中"
			}
			userStatus = append(userStatus, *v)
		}
		userStatus = privateutil.SortPrivateUserStatus(userStatus)
		lastWeekCount, _ = privateutil.GetStatusCount(userStatus)
		appAccountCount = privateutil.GetAppStatusCount(appStatus)
	}

	f.WriteData(http.StatusOK, api.StatusList{UserData: userStatus, AppData: appStatus, LastWeekCount: lastWeekCount, AccountCount: appAccountCount})
}

func GetTotalStatusHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[App Resource] Request to GetStatusHandler, %s", c.Request.URL.String())
	userResults, err := storage.MongoUS.GetUser(storage.UserFormat{Activate: "1"}, true, []string{}, []string{})
	if err != nil {
		f.WriteHTTPError(api.ErrInternalServer, err.Error())
		return
	}
	// get app status
	appResults, err := storage.MongoAS.GetApplist(storage.AppFormat{AccountID: ""})
	if err != nil {
		f.WriteHTTPError(api.ErrInternalServer, err.Error())
		return
	}

	res := api.PrivateTotalStatus{}
	res.AppIntegration = len(appResults)
	res.RegisterUser = len(userResults)
	res.Device, _ = deepstats.GetDeviceCount("all")
	{
		for _, v := range appResults {
			if v.AppID == "1652E90881C1FAE8" {
				continue
			}
			rawEvents := []string{"basic-event", "url-tem-overall-value", "sharelink-tem-overall-value"}
			convertedEvents := transmit.RequisiteAttrs(rawEvents, "t")
			aggregateResult, err := deepstats.GetChannelCounters(v.AppID, "all", convertedEvents, "t", "", "", "1", "")
			if err != nil {
				f.WriteHTTPError(api.ErrInternalServer, err.Error())
				return
			}
			eventRes := transmit.Calculate(rawEvents, convertedEvents, aggregateResult, "t", 1)
			res.LinkDemontration, _ = EventAddition(res.LinkDemontration, eventRes[0]["url-tem-overall-value"])
			res.LinkShare, _ = EventAddition(res.LinkShare, eventRes[0]["sharelink-tem-overall-value"])
			res.AppInstall, _ = EventAddition(res.AppInstall, eventRes[0]["match/install_with_params"])
			res.AppOpen, _ = EventAddition(res.AppOpen, eventRes[0]["match/open_with_params"])

		}
	}
	f.WriteData(http.StatusOK, res)
}
