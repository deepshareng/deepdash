package resource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/deepstats"
	"github.com/MISingularity/deepdash/pkg/errorutil"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/MISingularity/deepdash/pkg/sessionutil"
	"github.com/MISingularity/deepdash/pkg/transmit"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/gin"
)

var RoutesChannel = map[string]handler.Route{
	// GET Request
	"GetAppChannels": handler.Route{
		Name:        "GetAppChannels",
		Method:      "GET",
		Pattern:     "/apps/:appid/channels/statistics",
		HandlerFunc: GetAppChannelsHandler,
	},
	// from zhaohai
	// TODO:
	// Following handlers deal with the user channel resource
	// diferred from previous channel definition.
	// In order to avoid ambiguity, description of these handler should be changed.
	"PostNewAppChannel": handler.Route{
		Name:        "PostNewAppChannel",
		Method:      "POST",
		Pattern:     "/apps/:appid/channels",
		HandlerFunc: PostNewAppChannelHandler,
	},
	"DeleteAppChannel": handler.Route{
		Name:        "DeleteAppChannel",
		Method:      "DELETE",
		Pattern:     "/apps/:appid/channels/:channelname",
		HandlerFunc: DeleteAppChannelHandler,
	},
	"GetAppChannelInfo": handler.Route{
		Name:        "GetAppChannelInfo",
		Method:      "GET",
		Pattern:     "/apps/:appid/channels",
		HandlerFunc: GetAppChannelInfoHandler,
	},
	"PutChannelUrl": handler.Route{
		Name:        "PutChannelUrl",
		Method:      "PUT",
		Pattern:     "/apps/:appid/channelurl",
		HandlerFunc: PutChannelUrlHandler,
	},

	"GetChannelStatistics": handler.Route{
		Name:        "GetAppInfo",
		Method:      "GET",
		Pattern:     "/apps/:appid/statistics",
		HandlerFunc: GetChannelStatisticsHandler,
	},
}

var RoutesType = map[string]handler.Route{
	"GetType": handler.Route{
		Name:        "GetType",
		Method:      "GET",
		Pattern:     "/apps/:appid/types",
		HandlerFunc: TypesGetHandler,
	},
}

var RoutesSelectedItems = map[string]handler.Route{
	"GetSelectedItems": handler.Route{
		Name:        "GetSelectedItems",
		Method:      "GET",
		Pattern:     "/apps/:appid/selected_items",
		HandlerFunc: GetSelectedItemsHandler,
	},
	"PutSelectedItems": handler.Route{
		Name:        "PutSelectedItems",
		Method:      "PUT",
		Pattern:     "/apps/:appid/selected_items",
		HandlerFunc: PutSelectedItemsHandler,
	},
}

func DecodeChannelname(channelname string) (channelType string, channelName string, channelRemark string) {
	channelDetails := strings.Split(channelname, "_")
	channelType = "默认类型"
	channelName = "空"
	channelRemark = "空"

	if len(channelDetails) > 0 && channelDetails[0] != "" {
		channelType = channelDetails[0]
	}
	if len(channelDetails) > 1 && channelDetails[1] != "" {
		channelName = channelDetails[1]
	}
	if len(channelDetails) > 2 && channelDetails[2] != "" {
		channelRemark = channelDetails[2]
	}
	return
}

// GET Request
// URL : /apps/:appid/channels?event=install_with_params&event=open&...`
// If no of attr specify event, it will return all event counts of every channel.
// Currently, this API doesn't support query channels of given type, instead of returning all channels of an app. We will add this feature latter.
// timestamp: [start|end]: yyyy-mm-dd
func GetAppChannelsHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Channel Resource] Request to GetAppChannelsHandler, %s", c.Request.URL.String())
	log.Infof("[Channel Resource] Request Detail, AppID : %s", c.Param("appid"))
	appID := c.Param("appid")
	channels, err := deepstats.GetAppChannels(appID)
	if err != nil {
		f.WriteHTTPError(api.ErrDeepstatsGetAppChannelsFailed, err.Error())
		return
	}
	eventFilters := strings.Split(c.Query("event"), ",")
	if len(eventFilters) == 1 && eventFilters[0] == "" {
		eventFilters = []string{}
	}

	gran := "t"
	limit := 1
	start := c.Query("start")
	end := c.Query("end")
	log.Debugf("[Channel Resource] Request Detail, AppID: %s, gran: %s, start: %s, end: %s, limit: %s", c.Param("appid"), gran, start, end, limit)

	if start != "" && end != "" {
		gran = "d"
		// "2006-01-02" is required by package time, can not be changed!
		startTime, errStart := time.ParseInLocation("2006-01-02", start, time.Now().Location())
		endTime, errEnd := time.ParseInLocation("2006-01-02", end, time.Now().Location())
		if errStart != nil || errEnd != nil {
			log.Debugf("time error: %s %s", errStart, errEnd)
			f.WriteHTTPError(api.ErrDeepstatsGetCounterFailed, "invalid start|end")
			return
		}

		// end date in query in included, in deepstat end is not included
		endTime = endTime.Add(time.Hour * 24)

		start = strconv.FormatInt(startTime.Unix(), 10)
		end = strconv.FormatInt(endTime.Unix(), 10)
		limit = int(endTime.Sub(startTime).Hours() / 24)
	} else {
		start = ""
		end = ""
	}

	log.Debugf("[Channel Resource] Request Detail Converted, AppID: %s, gran: %s, start: %s, end: %s, limit: %s", c.Param("appid"), gran, start, end, limit)

	resdata := api.ChannelStatisticsGetReturnType{
		Data: []map[string]interface{}{},
	}
	convertEventFilters := transmit.RequisiteAttrs(eventFilters, "t")
	for _, v := range channels {
		// ignore overall channel
		if v.Channelname == "all" {
			continue
		}
		channelInfo := make(map[string]interface{})
		channelInfo["typename"], channelInfo["channelname"], channelInfo["remark"] = DecodeChannelname(v.Channelname)
		counters, err := deepstats.GetChannelCounters(appID, v.Channelname, convertEventFilters, gran, start, end, strconv.Itoa(limit), "")
		if err != nil {
			f.WriteHTTPError(api.ErrDeepstatsGetCounterFailed, fmt.Sprintf("[Channel Resource] Get deepstats channel counts failed! Err Msg=%v", err))
			return
		}
		res := transmit.Calculate(eventFilters, convertEventFilters, counters, gran, limit)

		// TODO: sum up
		for _, v := range res {
			for k, ch := range v {
				counterValue, err := strconv.ParseInt(ch, 10, 64)
				if err != nil {
					log.Errorf("[Channel Resource] Sum up statistics failed! Err Msg=%v", err)
					continue
				}
				if _, ok := channelInfo[k]; ok {
					channelInfo[k] = channelInfo[k].(int64) + counterValue
				} else {
					channelInfo[k] = counterValue
				}
			}
		}
		/*
			for k, ch := range res[0] {
				channelInfo[k] = ch
			}
		*/

		resdata.Data = append(resdata.Data, channelInfo)
	}
	log.Debugf("[--Result][handler][GetAppChannelsHandler] resdata=%v", resdata)
	f.WriteData(http.StatusOK, resdata)
}

// GET Request
// channel statistics
// URL : /apps/:appid/statistics
//   - parameters: "?appid=...&channel&gran=[d,w,y]&limit=10&event=install&event=...&start=y&end=z"
// Support several format rules requesting data, for pattern, we check every rule in turn.
// If there is more than one match, we will return the first one matched.
// 1. start=x&end=y
// 2. start=x&limit=y
// 3. end=x&limit=y
// 4. limit=x
func GetChannelStatisticsHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Channel Resource] Request to GetChannelStatisticsHandler, %s", c.Request.URL.String())
	eventFilters := strings.Split(c.Query("event"), ",")

	if len(eventFilters) == 1 && eventFilters[0] == "" {
		eventFilters = []string{}
	}

	log.Infof("[Channel Resource] Request Detail, AppID : %s, Channel : %s, Granularity : %s, Limit : %s, Start : %s, End : %s, Os : %s", c.Param("appid"), c.Query("channel"), c.Query("gran"), c.Query("limit"), c.Query("start"), c.Query("end"), c.Query("os"))
	limit, err := deepstats.GetLimit(c)
	if err != nil {
		f.WriteHTTPError(api.ErrDeepstatsGetLimitFailed, err.Error())
		return
	}
	convertEventFilters := transmit.RequisiteAttrs(eventFilters, c.Query("gran"))

	resdata := map[string][]map[string]string{}
	var oskeys []string

	if c.Query("groupby") == "os" {
		oskeys = []string{"", "ios", "android"}
	} else {
		oskeys = []string{""}
	}

	for _, osName := range oskeys {
		counters, err := deepstats.GetChannelCounters(c.Param("appid"), c.Query("channel"), convertEventFilters, c.Query("gran"), c.Query("start"), c.Query("end"), c.Query("limit"), osName)
		if err != nil {
			log.Errorf("[Channel Resource] Get deepstats channel counts failed! Err Msg=%v", err)
			f.WriteHTTPError(api.ErrDeepstatsGetCounterFailed, err.Error())
			return
		}
		// make key "all" in response
		if osName == "" {
			osName = "all"
		}
		resdata[osName] = transmit.Calculate(eventFilters, convertEventFilters, counters, c.Query("gran"), limit)
	}

	log.Debugf("[--Result][handler][GetChannelStatisticsHandler] resdata=%v", resdata)
	f.WriteData(http.StatusOK, resdata)
}

// GET /apps/:appid/types
func TypesGetHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	resdata := api.TypesGetRetrunType{Typelist: []api.Typeitem{api.Typeitem{TypeID: "1", TypeName: "默认类型"}}}
	f.WriteData(http.StatusOK, resdata)
}

// POST Request
// URL : /apps/:appid/channelinfo
func PostNewAppChannelHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Channel Resource] Request to PostNewAppChannelHandler, %s", c.Request.URL.String())

	appId := c.Param("appid")
	channelname := c.PostForm("channelname")
	channelurl := c.PostForm("channelurl")
	userId := sessionutil.GetAccountID(c)
	if userId == "" {
		f.WriteHTTPError(api.ErrBadRequestBody, "no accountid")
		return
	}
	log.Infof("[Channel Resource] Request Detail, AppID : %s, Channel : %s, ChannelUrl : %s", appId, channelname, channelurl)
	appchannels, err := storage.MongoCS.GetChannelList(storage.AppChannel{AppID: appId})
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	userIdBson := bson.ObjectIdHex(userId)
	userResults, err := storage.MongoUS.GetUser(storage.UserFormat{Id: userIdBson}, false, []string{}, []string{})
	if err != nil || len(userResults) < 1 {
		f.WriteHTTPError(api.ErrInternalServer, err.Error())
		return
	}

	userPermissionList := strings.Split(userResults[0].Permission, "_")
	unlimit := false
	for i := 0; i < len(userPermissionList); i++ {
		if userPermissionList[i] == "channelUnlimit" {
			unlimit = true
			break
		}
	}

	limitNum := 50
	if len(appchannels) >= limitNum && !unlimit {
		f.WriteData(http.StatusOK, map[string]string{"error": "推广链接数超过上限！扩充推广上限数目，请微信联系：huiyan258240"})
		return
	}
	err = storage.MongoCS.InsertChannel(storage.AppChannel{
		AppID:       appId,
		Channelname: channelname,
		Channelurl:  channelurl,
		CreateAt:    strconv.FormatInt(time.Now().Unix(), 10),
	})
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

// Delete Request
// URL : /apps/:appid/channelinfo
func DeleteAppChannelHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Channel Resource] Request to DeleteNewAppChannelHandler, %s", c.Request.URL.String())
	appId := c.Param("appid")
	channelname := c.Param("channelname")
	if channelname == "" || appId == "" {
		f.WriteHTTPError(api.ErrBadJSONBody, "")
		return
	}
	log.Infof("[Channel Resource] Request Detail, AppID : %s, Channel : %s", appId, channelname)
	err := deepstats.DeleteAppChannel(appId, channelname)
	if err != nil {
		f.WriteHTTPError(api.ErrDeepstatsConnectFailed, err.Error())
		return
	}
	err = storage.MongoCS.RemoveChannel(storage.AppChannel{AppID: appId, Channelname: channelname})
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

// GET Request
// URL : /apps/:appid/channelinfo
func GetAppChannelInfoHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Channel Resource] Request to GetAppChannelInfoHandler, %s", c.Request.URL.String())

	appId := c.Param("appid")
	log.Infof("[Channel Resource] Request Detail, AppID : %s", appId)

	appChannel, err := storage.MongoCS.GetChannelList(storage.AppChannel{AppID: appId})
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	resdata := api.AppChannelInfo{
		Data: appChannel,
	}
	log.Debugf("[--Result][handler][GetAppChannelInfoHandler] resdata=%v", resdata)
	f.WriteData(http.StatusOK, resdata)
}

// PUT Request
// URL : /apps/:appid/channelurl
func PutChannelUrlHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Channel Resource] Request to PutChannelUrlHandler, %s", c.Request.URL.String())

	appId := c.Param("appid")
	channelname := c.PostForm("channelname")
	channelurl := c.PostForm("channelurl")
	log.Infof("[Channel Resource] Request Detail, AppID : %s, ChannelName : %s, ChannelUrl : %s", appId, channelname, channelurl)

	if channelurl == "" || channelname == "" {
		return
	}

	appChannel, err := storage.MongoCS.GetChannel(storage.AppChannel{AppID: appId, Channelname: channelname})

	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	// insert new appchannel
	upsertApp := storage.AppChannel{
		AppID:       appChannel.AppID,
		Channelname: appChannel.Channelname,
		Channelurl:  channelurl,
	}

	upsertinfo, err := json.Marshal(upsertApp)
	if err != nil {
		log.Errorf("[Channel Resource] Marshal upsert info failed! Err Msg=%v", err)
		return
	}
	log.Infof("[Channel Resource] Request Detail, UpsertInfo : %s", string(upsertinfo))

	err = storage.MongoCS.UpdateChannel(storage.AppChannel{AppID: appId}, upsertApp, true)
	httperror = errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}

	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}

// Get Request
// URL : /apps/:appid/selected_items

func GetSelectedItemsHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Items Resource] Request to GetSelectedItemsHandler, %s", c.Request.URL.String())
	search := storage.AppSelectedItems{AppID: c.Param("appid"), AccountID: sessionutil.GetAccountID(c)}
	log.Infof("[Items Resource] Request Detail, accountid=%s, appid=%s", sessionutil.GetAccountID(c), c.Param("appid"))
	res, err := storage.MongoASS.GetAppSelectedItems(search)
	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, res)
}

// Put Request
// URL : /apps/:appid/selected_items

func PutSelectedItemsHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Items Resource] Request to PutSelectedItemsHandler, %s", c.Request.URL.String())
	var items storage.AppSelectedItems

	err := c.BindJSON(&items)
	if err != nil {
		f.WriteHTTPError(api.ErrBadJSONBody, err.Error())
		return
	}
	items.AccountID = sessionutil.GetAccountID(c)
	items.AppID = c.Param("appid")
	search := storage.AppSelectedItems{AppID: c.Param("appid"), AccountID: sessionutil.GetAccountID(c)}
	log.Infof("[Items Resource] Request Detail, upsert=%v", items)
	err = storage.MongoASS.UpdateAppSelectedItems(search, items, true)

	httperror := errorutil.ProcessMongoError(f, err)
	if httperror {
		return
	}
	f.WriteData(http.StatusOK, api.SimpleSuccessResponse{Success: true})
}
