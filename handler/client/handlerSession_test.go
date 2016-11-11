package client

import (
	"net/http"
	"path"
	"testing"

	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/testutil"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestGetSessionUserInfo(t *testing.T) {
	dbName := "test-getsessionuserinfo"
	collName := "getsessionuserinfo"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()
	storage.MongoUS, _ = storage.NewMongoUserStorageService(c)
	storage.PrepareUserItems(storage.MongoUS)

	testcases := []struct {
		requestPath string
		username    string

		wcode int
		wbody string
	}{
		{
			"/session/getuser",
			"test1",

			http.StatusOK,
			`{"username":"test1"}` + "\n",
		},
		{
			"/session/getuser",
			"test9",

			http.StatusOK,
			`{"username":""}` + "\n",
		},
	}

	var store = sessions.NewCookieStore([]byte("secret"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	router.Use(testutil.AccountIDDebugSetting)
	testroute := RoutesSession["GetSessionUserInfo"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	//router.Use(gin.ErrorLogger())

	for i, tt := range testcases {
		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		results, _ := storage.MongoUS.GetUser(storage.UserFormat{Username: tt.username}, false, []string{}, []string{})
		u := storage.UserFormat{}
		if len(results) > 0 {
			u = results[0]
		}
		handler.Displayname = u.Username
		w := testutil.HandleWithBody(router, testroute.Method, url, "")
		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
		}

		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
	}
}
