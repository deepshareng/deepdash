package resource

import (
	"encoding/json"
	"net/http"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/MISingularity/deepdash/api"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/deepstats"
	"github.com/MISingularity/deepdash/pkg/testutil"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestGetAppChannelsHandler(t *testing.T) {
	dbName := "test-getappchannelshandler"
	collName := "getappchannelshandler"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	storage.MongoAS = storage.NewMongoAppStorageService(c)
	storage.PrepareAppItems(storage.MongoAS)

	testcases := []struct {
		requestPath string
		appid       string
		body        string

		mockAppInfoResponseCode []int
		mockAppInfoResponseBody []string

		wcode int
		res   []map[string]interface{}
	}{
		// update ios info
		{
			"/apps/appid1/channels/statistics",
			"appid1",
			"",

			[]int{http.StatusOK, http.StatusOK},
			[]string{`{"AppID":"appid1","Channels":["new_x_y"]}`, `{"Counters":[{"Event":"match/install_with_params","Counts":[{"Count":1}]}]}`},

			http.StatusOK,
			[]map[string]interface{}{map[string]interface{}{"typename": "new", "channelname": "x", "remark": "y", "match/install_with_params": float64(1)}},
		},

		{
			"/apps/appid3/channels/statistics",
			"appid3",
			"",

			[]int{http.StatusOK, http.StatusOK},
			[]string{`{"AppID":"appid3","Channels":["new"]}`, `{"Counters":[{"Event":"match/install_with_params","Counts":[{"Count":1}]},{"Event":"match/open_with_params","Counts":[{"Count":1}]}]}`},

			http.StatusOK,
			[]map[string]interface{}{map[string]interface{}{"typename": "new", "channelname": "空", "remark": "空", "match/install_with_params": float64(1), "match/open_with_params": float64(1)}},
		},

		{
			"/apps/appid10/channels/statistics",
			"appid10",
			"",

			[]int{http.StatusOK, http.StatusOK},
			[]string{`{"AppID":"appid3","Channels":["new_1_2"]}`, `{"Counters":[]}`},

			http.StatusOK,
			[]map[string]interface{}{map[string]interface{}{"typename": "new", "channelname": "1", "remark": "2"}},
		},
	}
	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	testroute := RoutesChannel["GetAppChannels"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	//router.Use(gin.ErrorLogger())

	for i, tt := range testcases {
		serverMock, _, _ := testutil.MockResponse(tt.mockAppInfoResponseCode, tt.mockAppInfoResponseBody)
		defer serverMock.Close()
		deepstats.DEEPSTATSD_URL = serverMock.URL

		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		w := testutil.HandleWithBody(router, testroute.Method, url, tt.body)

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
			continue
		}
		var resdata api.ChannelStatisticsGetReturnType
		err := json.Unmarshal(w.Body.Bytes(), &resdata)
		if err != nil {
			t.Errorf("#%d: Unmarshal response data failed! Err Msg=%v", i, err)
		}
		if len(tt.res) != len(resdata.Data) {
			t.Errorf("#%d: HTTP response array length mismatch, get = %v, want = %v", i, resdata, tt.res)
			continue
		}
		for j := range tt.res {
			if len(tt.res[j]) != len(resdata.Data[j]) {
				t.Errorf("#%d: Counters result array #%d map length mismatch, get = %v, want = %v", i, j, resdata, tt.res)
				continue
			}
			for k, v := range tt.res[j] {
				if v != resdata.Data[j][k] {
					t.Errorf("#%d: Counters result array #%d find key %s value mismatch! Get value = %s. Want value = %s", i, j, k, v, resdata.Data[j][k])
				}
			}
		}

	}
}

func TestDeleteAppChannelHandler(t *testing.T) {
	dbName := "test-deleteappchannelhandler"
	collName := "deleteappchannelhandler"
	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	storage.MongoCS = storage.NewMongoAppChannelService(c)
	storage.PrepareChannelItems(storage.MongoCS)
	testcases := []struct {
		requestPath string
		appid       string
		channel     string

		mockDeepstatsResponseCode []int
		mockDeepstatsResponseBody []string

		wcode int
	}{
		{
			"/apps/appid1/channels/channelname1",
			"appid1",
			"channelname1",

			[]int{http.StatusOK},
			[]string{""},

			http.StatusOK,
		},
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	testroute := RoutesChannel["DeleteAppChannel"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	for i, tt := range testcases {
		serverMock, _, _ := testutil.MockResponse(tt.mockDeepstatsResponseCode, tt.mockDeepstatsResponseBody)
		defer serverMock.Close()
		deepstats.DEEPSTATSD_URL = serverMock.URL

		var requesturl string

		requesturl = "http://" + path.Join("exmaple.com", tt.requestPath)
		w := testutil.HandleWithBody(router, testroute.Method, requesturl, "")

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
			continue
		}

		_, err := storage.MongoCS.GetChannel(storage.AppChannel{
			AppID:       tt.appid,
			Channelname: tt.channel,
		})
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Errorf("#%d: Delete failed", i)
			continue
		}
	}

}

func TestGetChannelStatisticsHandler(t *testing.T) {
	dbName := "test-getchannelinfohandler"
	collName := "getchannelinfohandler"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	storage.MongoAS = storage.NewMongoAppStorageService(c)
	storage.PrepareAppItems(storage.MongoAS)

	testcases := []struct {
		requestPath string
		params      string
		appid       string
		body        string

		mockAppInfoResponseCode []int
		mockAppInfoResponseBody []string

		wcode int
		//res   []map[string]string
		res map[string][]map[string]string
	}{
		// update ios info
		{
			"/apps/appid/statistics",
			"?appid=appid&gran=d&limit=2&channel=channel1&groupby=os",
			"appid1",
			"",

			[]int{http.StatusOK, http.StatusOK, http.StatusOK},
			[]string{
				`{"Counters":[{"Event":"match/install_with_params","Counts":[{"Count":1}, {"Count":0}]}]}`,
				`{"Counters":[{"Event":"match/install_with_params","Counts":[{"Count":1}, {"Count":0}]}]}`,
				`{"Counters":[{"Event":"match/install_with_params","Counts":[{"Count":1}, {"Count":0}]}]}`,
			},

			http.StatusOK,
			//[]map[string]string{map[string]string{"match/install_with_params": "0"}, map[string]string{"match/install_with_params": "1"}},
			map[string][]map[string]string{
				"all":     []map[string]string{map[string]string{"match/install_with_params": "0"}, map[string]string{"match/install_with_params": "1"}},
				"android": []map[string]string{map[string]string{"match/install_with_params": "0"}, map[string]string{"match/install_with_params": "1"}},
				"ios":     []map[string]string{map[string]string{"match/install_with_params": "0"}, map[string]string{"match/install_with_params": "1"}},
			},
		},

		{
			"/apps/appid/statistics",
			"?appid=appid1&gran=d&limit=3&channel=channel1",
			"appid1",
			"",
			[]int{http.StatusOK},
			[]string{
				`{"Counters":[{"Event":"match/install_with_params","Counts":[{"Count":1},{"Count":2},{"Count":3}]},{"Event":"match/install","Counts":[{"Count":11},{"Count":12},{"Count":13}]},{"Event":"match/open_with_params","Counts":[{"Count":21},{"Count":22},{"Count":23}]},{"Event":"match/open","Counts":[{"Count":31},{"Count":32},{"Count":33}]}]}`,
			},

			http.StatusOK,
			map[string][]map[string]string{
				"all": []map[string]string{
					map[string]string{"match/open": "33", "match/open_with_params": "23", "match/install": "13", "match/install_with_params": "3"},
					map[string]string{"match/open": "32", "match/open_with_params": "22", "match/install": "12", "match/install_with_params": "2"},
					map[string]string{"match/open": "31", "match/open_with_params": "21", "match/install": "11", "match/install_with_params": "1"},
				},
			},
		},
		{
			"/apps/appid/statistics",
			"?appid=appid1&gran=d&limit=3&channel=channel1&groupby=os",
			"appid1",
			"",

			[]int{http.StatusOK, http.StatusOK, http.StatusOK},
			[]string{
				`{"Counters":[]}`,
				`{"Counters":[]}`,
				`{"Counters":[]}`,
			},

			http.StatusOK,
			//[]map[string]string{{}, {}, {}},
			map[string][]map[string]string{
				"all":     []map[string]string{{}, {}, {}},
				"android": []map[string]string{{}, {}, {}},
				"ios":     []map[string]string{{}, {}, {}},
			},
		},
	}

	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	testroute := RoutesChannel["GetChannelStatistics"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	//router.Use(gin.ErrorLogger())

	for i, tt := range testcases {
		serverMock, _, _ := testutil.MockResponse(tt.mockAppInfoResponseCode, tt.mockAppInfoResponseBody)
		defer serverMock.Close()
		deepstats.DEEPSTATSD_URL = serverMock.URL

		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath) + tt.params
		w := testutil.HandleWithBody(router, testroute.Method, url, tt.body)

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
			continue
		}

		//resdata := api.ChannelGetReturnType{}
		resdata := map[string][]map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &resdata)
		if err != nil {
			t.Errorf("#%d: Unmarshal response data failed! Err Msg=%v", i, err)
		}

		for o, l := range tt.res {
			if len(l) != len(resdata[o]) {
				t.Errorf("#%d: HTTP response counters result array length mismatch, get = %v, want = %v", i, resdata, tt.res)
				continue
			}
			for j := range l {
				if len(l[j]) != len(resdata[o][j]) {
					t.Errorf("#%d: Counters result array #%d map length mismatch, get = %v, want = %v", i, j, resdata, tt.res)
					continue
				}
				for k, v := range l[j] {
					if v != resdata[o][j][k] {
						t.Errorf("#%d: Counters result key %s array #%d map find key %s value mismatch! Get value = %s. Want value = %s", i, o, j, k, v, resdata[o][j][k])
					}
				}
			}

		}
	}
}

func TestGetSelectedItemsHandler(t *testing.T) {
	dbName := "test-getselecteditemshandler"
	collName := "geselecteditemshandler"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	storage.MongoASS = storage.NewMongoAppSelectedItemsService(c)
	storage.PrepareAppSelectedItems(storage.MongoASS)
	testcases := []struct {
		requestPath string
		accountid   string

		wcode int
		wbody string
	}{
		{
			"/apps/appid1/selected_items",
			"acc1",

			http.StatusOK,
			`{"appid":"appid1","accountid":"acc1","events":["a1","b1","c1"],"displays":["da1","db1","dc1"]}` + "\n",
		},
		{
			"/apps/appid2/selected_items",
			"acc1",

			http.StatusOK,
			`{"appid":"","accountid":"","events":null,"displays":null}` + "\n",
		},
	}

	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	router.Use(testutil.AccountIDDebugSetting)
	testroute := RoutesSelectedItems["GetSelectedItems"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	for i, tt := range testcases {
		handler.AccountID = tt.accountid
		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		w := testutil.HandleWithBody(router, testroute.Method, url, "")
		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}

		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
	}
}

func TestPutSelectedItemsHandler(t *testing.T) {
	dbName := "test-putselecteditemshandler"
	collName := "putselecteditemshandler"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	storage.MongoASS = storage.NewMongoAppSelectedItemsService(c)
	storage.PrepareAppSelectedItems(storage.MongoASS)
	testcases := []struct {
		requestPath string
		accountid   string
		appid       string
		header      map[string]string
		body        string

		wcode   int
		wbody   string
		updated storage.AppSelectedItems
	}{
		{
			"/apps/appid1/selected_items",
			"acc1",
			"appid1",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			`{"events":["a","b"],"displays":["aaa","b"]}`,

			http.StatusOK,
			`{"success":true}` + "\n",
			storage.AppSelectedItems{
				AppID:     "appid1",
				AccountID: "acc1",
				Events:    []string{"a", "b"},
				Displays:  []string{"aaa", "b"},
			},
		},
		{
			"/apps/appid4/selected_items",
			"acc2",
			"appid4",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			`{"events":["aaaa","bbbb","ccc"],"displays":["aaa","b","ccc"]}`,

			http.StatusOK,
			`{"success":true}` + "\n",
			storage.AppSelectedItems{
				AppID:     "appid4",
				AccountID: "acc2",
				Events:    []string{"aaaa", "bbbb", "ccc"},
				Displays:  []string{"aaa", "b", "ccc"},
			},
		},
	}

	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	router.Use(testutil.AccountIDDebugSetting)
	testroute := RoutesSelectedItems["PutSelectedItems"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	//router.Use(gin.ErrorLogger())

	for i, tt := range testcases {
		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		handler.AccountID = tt.accountid
		w := testutil.HandleWithRequestInfo(router, testroute.Method, url, tt.body, tt.header, "ip1:port1")
		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}

		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
		search := storage.AppSelectedItems{AppID: tt.appid, AccountID: tt.accountid}
		res, err := storage.MongoASS.GetAppSelectedItems(search)
		if err != nil {
			t.Errorf("#%d: Get Selected Items Failed! Err Msg = %v", i, err)
			continue
		}
		if !reflect.DeepEqual(res, tt.updated) {
			t.Errorf("#%d : App selected items storage update failed, want = %v, get = %v", i, tt.updated, res)
		}

	}
}
