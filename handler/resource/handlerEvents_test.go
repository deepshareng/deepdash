package resource

import (
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/MISingularity/deepdash/pkg/deepstats"
	"github.com/MISingularity/deepdash/pkg/testutil"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestGetAppEventsHandler(t *testing.T) {
	testcases := []struct {
		requestPath string
		appid       string
		body        string

		mockAppInfoResponseCode []int
		mockAppInfoResponseBody []string

		wcode int
		wbody string
	}{
		// update ios info
		{
			"/apps/appid1/events",
			"appid1",
			"",

			[]int{http.StatusOK, http.StatusOK},
			[]string{`{"AppID":"appid1","Events":["/v2/counters/install","/v2/counters/open"]}`},

			http.StatusOK,
			`{"eventlist":[{"event":"match/install_with_params","display":"新用户下载"},{"event":"match/open_with_params","display":"老用户打开"},{"event":"3-day-retention","display":"三日留存"},{"event":"7-day-retention","display":"七日留存"},{"event":"/v2/counters/install","display":"install"},{"event":"/v2/counters/open","display":"open"}]}` + "\n",
		},
	}
	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	testroute := RoutesEvent["GetAppEvents"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	//router.Use(gin.ErrorLogger())

	for i, tt := range testcases {
		serverMock, _, _ := testutil.MockResponse(tt.mockAppInfoResponseCode, tt.mockAppInfoResponseBody)
		defer serverMock.Close()
		deepstats.DEEPSTATSD_URL = serverMock.URL

		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		fmt.Println(url)
		w := testutil.HandleWithBody(router, testroute.Method, url, tt.body)

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
			continue
		}
		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}

	}
}
