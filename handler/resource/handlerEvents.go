package resource

import (
	"net/http"
	"strings"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/deepstats"
	"github.com/MISingularity/deepdash/pkg/httputil"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/gin-gonic/gin"
)

var RoutesEvent = map[string]handler.Route{
	// GET Request
	"GetAppEvents": handler.Route{
		Name:        "GetAppEvents",
		Method:      "GET",
		Pattern:     "/apps/:appid/events",
		HandlerFunc: GetAppEventsHandler,
	},
}

// GET Request
// URL : /apps/:appid/events
func GetAppEventsHandler(c *gin.Context) {
	f := httputil.NewGinframework(c)
	log.Infof("[Event Resource] Request to GetAppChannelsHandler, %s", c.Request.URL.String())
	log.Infof("[Event Resource] Request Detail, AppID : %s", c.Param("appid"))
	appID := c.Param("appid")
	events, err := deepstats.GetAppEvents(appID)
	if err != nil {
		f.WriteHTTPError(api.ErrDeepstatsConnectFailed, err.Error())
		return
	}
	resdata := api.Eventlist{
		Eventlist: []api.Event{
			api.Event{
				Event:   "match/install_with_params",
				Display: "新用户下载",
			},
			api.Event{
				Event:   "match/open_with_params",
				Display: "老用户打开",
			},
			api.Event{
				Event:   "3-day-retention",
				Display: "三日留存",
			},
			api.Event{
				Event:   "7-day-retention",
				Display: "七日留存",
			},
		},
	}
	for _, v := range events.Events {
		display := v
		if strings.HasPrefix(display, "/v2/counters/") {
			display = display[13:]
		}
		resdata.Eventlist = append(resdata.Eventlist, api.Event{Event: v, Display: display})
	}
	log.Debugf("[--Result][handler][GetAppEventsHandler] resdata=%v", resdata)
	f.WriteData(http.StatusOK, resdata)
}
