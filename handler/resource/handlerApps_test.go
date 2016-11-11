package resource

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"path"
	"reflect"
	"testing"

	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/deepstats"
	"github.com/MISingularity/deepdash/pkg/testutil"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestGetApplistHandler(t *testing.T) {
	dbName := "test-getapplishandler"
	collName := "getapplishandler"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	storage.MongoAS = storage.NewMongoAppStorageService(c)
	storage.PrepareAppItems(storage.MongoAS)

	testcases := []struct {
		requestPath string
		account     string
		body        string

		mockAppInfoResponseCode []int
		mockAppInfoResponseBody []string

		wcode int
		wbody string
	}{
		{
			"/apps",
			"account1",
			"",

			[]int{http.StatusOK, http.StatusOK},
			[]string{`{"AppID":"appid1","Channels":["new"]}`, `{"AppID":"appid1","Channels":["new1","new2"]}`},

			http.StatusOK,
			`{"applist":[{"appid":"appid1","appname":"app1","iconurl":"","channelinfo":null},{"appid":"appid2","appname":"app2","iconurl":"","channelinfo":null}]}` + "\n",
		},
		{
			"/apps",
			"account3",
			"",

			[]int{http.StatusOK},
			[]string{`{"AppID":"appid1","Channels":["new"]}`},

			http.StatusOK,
			`{"applist":[{"appid":"appid3","appname":"app3","iconurl":"","channelinfo":null}]}` + "\n",
		},
		{
			"/apps",
			"account2",
			"",

			[]int{http.StatusOK},
			[]string{``},

			http.StatusOK,
			`{"applist":[]}` + "\n",
		},
	}
	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	router.Use(testutil.AccountIDDebugSetting)
	testroute := RoutesApp["GetApplist"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	for i, tt := range testcases {
		serverMock, _, _ := testutil.MockResponse(tt.mockAppInfoResponseCode, tt.mockAppInfoResponseBody)
		defer serverMock.Close()
		deepstats.DEEPSTATSD_URL = serverMock.URL

		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		handler.AccountID = tt.account
		w := testutil.HandleWithBody(router, testroute.Method, url, tt.body)
		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}

		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
	}
}

func TestGetAppinfoHandler(t *testing.T) {
	dbName := "test-getappinfohandler"
	collName := "getappinfohandler"

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

		wcode int
		wbody string
	}{
		{
			"/apps",
			"appid1",
			"",

			http.StatusOK,
			`{"appID":"appid1","shortID":"","accountID":"account1","appName":"app1","pkgName":"pkg1","iosBundler":"","iosScheme":"","iosDownloadUrl":"","iosUniversalLinkEnable":"","iosTeamID":"","iosYYBEnableBelow9":"","iosYYBEnableAbove9":"","androidPkgname":"","androidScheme":"","androidHost":"","androidAppLink":"","androidDownloadUrl":"","androidIsDownloadDirectly":"","androidYYBEnable":"","androidSHA256":"","yyburl":"","yybenable":"","attributionPushUrl":"","iconUrl":"","downloadTitle":"","downloadMsg":"","theme":"","userConfBgWeChatAndroidTipUrl":"","userConfBgWeChatIosTipUrl":"","forceDownload":""}` + "\n",
		},
		{
			"/apps",
			"appid2",
			"",

			http.StatusOK,
			`{"appID":"appid2","shortID":"","accountID":"account1","appName":"app2","pkgName":"pkg2","iosBundler":"","iosScheme":"","iosDownloadUrl":"","iosUniversalLinkEnable":"","iosTeamID":"","iosYYBEnableBelow9":"","iosYYBEnableAbove9":"","androidPkgname":"","androidScheme":"","androidHost":"","androidAppLink":"","androidDownloadUrl":"","androidIsDownloadDirectly":"","androidYYBEnable":"","androidSHA256":"","yyburl":"","yybenable":"","attributionPushUrl":"","iconUrl":"","downloadTitle":"","downloadMsg":"","theme":"","userConfBgWeChatAndroidTipUrl":"","userConfBgWeChatIosTipUrl":"","forceDownload":""}` + "\n",
		},
		{
			"/apps",
			"appid10",
			"",

			http.StatusOK,
			`{"appID":"","shortID":"","accountID":"","appName":"","pkgName":"","iosBundler":"","iosScheme":"","iosDownloadUrl":"","iosUniversalLinkEnable":"","iosTeamID":"","iosYYBEnableBelow9":"","iosYYBEnableAbove9":"","androidPkgname":"","androidScheme":"","androidHost":"","androidAppLink":"","androidDownloadUrl":"","androidIsDownloadDirectly":"","androidYYBEnable":"","androidSHA256":"","yyburl":"","yybenable":"","attributionPushUrl":"","iconUrl":"","downloadTitle":"","downloadMsg":"","theme":"","userConfBgWeChatAndroidTipUrl":"","userConfBgWeChatIosTipUrl":"","forceDownload":""}` + "\n",
		},
	}

	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	testroute := RoutesApp["GetAppinfo"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	//router.Use(gin.ErrorLogger())

	for i, tt := range testcases {
		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath, tt.appid)

		w := testutil.HandleWithBody(router, testroute.Method, url, tt.body)
		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}

		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
	}
}

func TestPutAppHandler(t *testing.T) {
	dbName := "test-putapphandler"
	collName := "putapphandler"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	cSeq := session.DB(dbName).C("seq")
	cSeq.Insert(bson.M{"seq_key": "app", "seq_val": 0})
	defer session.DB(dbName).DropDatabase()

	storage.MongoAS = storage.NewMongoAppStorageService(c)
	storage.MongoSS = storage.NewMongoSequenceService(cSeq)
	storage.PrepareAppItems(storage.MongoAS)

	testcases := []struct {
		requestPath string
		appid       string
		header      map[string]string
		body        string
		remoteAddr  string

		mockAppInfoResponseCode []int
		mockAppInfoResponseBody []string

		wcode     int
		wbody     string
		updateApp storage.AppFormat
	}{
		// update ios info
		{
			"/apps",
			"appid1",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			"appName=app1&pkgName=pkg2&iosBundler=bundler&iosScheme=scheme&iosDownloadUrl=url&iosUniversalLinkEnable=true&iosTeamID=1&androidPkgname=&androidScheme=&androidHost=&androidAppLink=&androidDownloadUrl=&androidIsDownloadDirectly=false&yyburl=&yybenable=true&iosYYBEnableBelow9=true&iosYYBEnableAbove9=true",
			"ip:port1",

			[]int{http.StatusOK},
			[]string{`{"AppID":"7713337217A6E150","Android":{"Scheme":"deepshare","Host":"com.singulariti.deepsharedemo","Pkg":"com.singulariti.deepsharedemo","DownloadUrl":"","IsDownloadDirectly":"true"},"Ios":{"Scheme":"deepsharedemo","DownloadUrl":""},"YYBUrl":"","YYBEnable":"false","forceDownload":"false"}`},

			http.StatusOK,
			`{"success":true}` + "\n",
			storage.AppFormat{
				AppID:     "appid1",
				ShortID:   "aaab",
				AccountID: "account1",
				AppName:   "app1",
				PkgName:   "pkg2",

				AndroidScheme:             "dsappid1",
				AndroidAppLink:            "false",
				AndroidIsDownloadDirectly: "false",
				AndroidDownloadUrl:        "",
				AndroidYYBEnable:          "false",

				IosBundler:             "bundler",
				IosScheme:              "dsappid1",
				IosDownloadUrl:         "url",
				IosUniversalLinkEnable: "true",
				IosTeamID:              "1",
				IosYYBEnableBelow9:     "true",
				IosYYBEnableAbove9:     "true",

				YYBenable:     "true",
				ForceDownload: "false",
			},
		},

		// update android info
		{
			"/apps",
			"appid3",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			"appName=appw&pkgName=pkgchange&iosBundler=&iosScheme=&iosDownloadUrl=&iosUniversalLinkEnable=&iosTeamID=&androidPkgname=pkg&androidScheme=aschema&androidHost=host&androidSHA256=123&androidDownloadUrl=testURL&androidIsDownloadDirectly=true&yyburl=yyburl&yybenable=true&androidYYBEnable=true",
			"ip:port1",

			[]int{http.StatusOK},
			[]string{`{"AppID":"7713337217A6E150","Android":{"Scheme":"deepshare","Host":"com.singulariti.deepsharedemo","Pkg":"com.singulariti.deepsharedemo","DownloadUrl":""},"Ios":{"Scheme":"deepsharedemo","DownloadUrl":""},"YYBUrl":"","YYBEnable":"false","ForceDownload":"false"}`},

			http.StatusOK,
			`{"success":true}` + "\n",
			storage.AppFormat{
				AppID:     "appid3",
				ShortID:   "aaac",
				AccountID: "account3",
				AppName:   "appw",
				PkgName:   "pkgchange",

				IosScheme:              "dsappid3",
				IosUniversalLinkEnable: "false",
				IosYYBEnableBelow9:     "false",
				IosYYBEnableAbove9:     "false",

				AndroidPkgname:            "pkg",
				AndroidScheme:             "dsappid3",
				AndroidHost:               "host",
				AndroidAppLink:            "true",
				AndroidSHA256:             "123",
				AndroidDownloadUrl:        "testURL",
				AndroidIsDownloadDirectly: "true",
				AndroidYYBEnable:          "true",
				YYBenable:                 "true",
				YYBurl:                    "yyburl",
				ForceDownload:             "false",
			},
		},
		// create new one
		{
			"/apps",
			"appidnew",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			"appName=appw&pkgName=pkgchange&apptype=&iosBundler=bundler&iosScheme=scheme&iosDownloadUrl=url&iosUniversalLinkEnable=true&iosTeamID=1&androidPkgname=pkg&androidScheme=aschema&androidHost=host&androidDownloadUrl=testURL&androidIsDownloadDirectly=true&yyburl=yyburl&yybenable=true&theme=122&userConfBgWeChatAndroidTipUrl=123&userConfBgWeChatIosTipUrl=124",
			"ip:port1",

			[]int{http.StatusOK},
			[]string{`{"AppID":"7713337217A6E150","Android":{"Scheme":"deepshare","Host":"com.singulariti.deepsharedemo","Pkg":"com.singulariti.deepsharedemo","DownloadUrl":""},"Ios":{"Scheme":"deepsharedemo","DownloadUrl":""},"YYBUrl":"","YYBEnable":false}`},

			http.StatusOK,
			`{"success":true}` + "\n",
			storage.AppFormat{},
		},
	}
	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	router.Use(testutil.AccountIDDebugSetting)
	testroute := RoutesApp["PutApp"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	for i, tt := range testcases {
		serverMock, _, _ := testutil.MockResponse(tt.mockAppInfoResponseCode, tt.mockAppInfoResponseBody)
		defer serverMock.Close()
		storage.RemoteAppInfoAS = storage.NewRemoteAppInfoService(serverMock.URL)

		handler.AccountID = ""
		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath, tt.appid)
		w := testutil.HandleWithRequestInfo(router, testroute.Method, url, tt.body, tt.header, tt.remoteAddr)

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
			continue
		}
		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
		result, _ := storage.MongoAS.GetApp(storage.AppFormat{AppID: tt.appid})

		if !reflect.DeepEqual(result, tt.updateApp) {
			t.Errorf("#%d : App storage update failed, want = %v, get = %v", i, tt.updateApp, result)
		}
	}
}

func TestPostAppHandler(t *testing.T) {
	dbName := "test-postapphandler"
	collName := "postapphandler"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	cSeq := session.DB(dbName).C("seq")
	cSeq.Insert(bson.M{"seq_key": "app", "seq_val": 0})
	defer session.DB(dbName).DropDatabase()
	storage.MongoAS = storage.NewMongoAppStorageService(c)
	storage.MongoSS = storage.NewMongoSequenceService(cSeq)
	storage.PrepareAppItems(storage.MongoAS)

	testcases := []struct {
		requestPath string
		header      map[string]string
		body        string
		remoteAddr  string

		accountID string

		mockAppInfoResponseCode []int
		mockAppInfoResponseBody []string

		wcode     int
		updateApp storage.AppFormat
	}{
		// create new one
		{
			"/apps",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			"appName=appw&pkgName=pkgchange&apptype=&iosBundler=bundler&iosScheme=scheme&iosDownloadUrl=url&iosUniversalLinkEnable=true&iosTeamID=1&androidPkgname=pkg&androidScheme=aschema&androidHost=host&androidDownloadUrl=testURL&androidIsDownloadDirectly=true&yyburl=yyburl&yybenable=true&theme=122&userConfBgWeChatAndroidTipUrl=123&userConfBgWeChatIosTipUrl=124",
			"ip:port1",

			"account10",

			[]int{http.StatusOK},
			[]string{`{"AppID":"7713337217A6E150","Android":{"Scheme":"deepshare","Host":"com.singulariti.deepsharedemo","Pkg":"com.singulariti.deepsharedemo","DownloadUrl":""},"Ios":{"Scheme":"deepsharedemo","DownloadUrl":""},"YYBUrl":"","YYBEnable":false}`},

			http.StatusOK,
			storage.AppFormat{
				AppID:     "appidnew",
				ShortID:   "aaab",
				AccountID: "account10",
				AppName:   "appw",
				PkgName:   "pkgchange",
				Theme:     "122",
				UserConfBgWeChatAndroidTipUrl: "123",
				UserConfBgWeChatIosTipUrl:     "124",

				IosBundler:             "bundler",
				IosScheme:              "dsappidnew",
				IosDownloadUrl:         "url",
				IosUniversalLinkEnable: "true",
				IosTeamID:              "1",
				IosYYBEnableBelow9:     "false",
				IosYYBEnableAbove9:     "false",

				AndroidPkgname:            "pkg",
				AndroidScheme:             "dsappidnew",
				AndroidHost:               "host",
				AndroidAppLink:            "false",
				AndroidDownloadUrl:        "testURL",
				AndroidIsDownloadDirectly: "true",
				AndroidYYBEnable:          "false",

				YYBenable: "true",
				YYBurl:    "yyburl",
			},
		},
	}
	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	router.Use(testutil.AccountIDDebugSetting)
	testroute := RoutesApp["PostApp"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	for i, tt := range testcases {
		serverMock, _, _ := testutil.MockResponse(tt.mockAppInfoResponseCode, tt.mockAppInfoResponseBody)
		defer serverMock.Close()
		storage.RemoteAppInfoAS = storage.NewRemoteAppInfoService(serverMock.URL)

		handler.AccountID = tt.accountID
		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		w := testutil.HandleWithRequestInfo(router, testroute.Method, url, tt.body, tt.header, tt.remoteAddr)

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
			continue
		}
		res := map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &res)
		if err != nil {
			t.Errorf("#%d: Unmarshl failed! Body = %v, err = %v", i, string(w.Body.Bytes()), err)
			return
		}
		tt.updateApp.AppID = res["appid"]
		tt.updateApp.AndroidScheme = "ds" + res["appid"]
		tt.updateApp.IosScheme = "ds" + res["appid"]
		if w.Code != http.StatusOK {
			continue
		}
		result, err := storage.MongoAS.GetApp(storage.AppFormat{AppID: res["appid"]})
		if err != nil {
			t.Errorf("#%d: Get App Failed! Err Msg = %v", i, err)
			continue
		}
		if !reflect.DeepEqual(result, tt.updateApp) {
			t.Errorf("#%d : App storage update failed, want = %v, get = %v", i, tt.updateApp, result)
		}
	}
}
