package client

import (
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/MISingularity/deepdash/pkg/testutil"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestResetPasswordHandler(t *testing.T) {
	dbName := "test-resetpassword-handler"
	session := testutil.MustNewLocalSession()
	defer session.Close()
	tokenColl := session.DB(dbName).C("token")
	userColl := session.DB(dbName).C("user")
	defer session.DB(dbName).DropDatabase()
	storage.MongoUS, _ = storage.NewMongoUserStorageService(userColl)
	storage.MongoTS, _ = storage.NewMongoTokenStorageService(tokenColl, time.Duration(0))
	storage.MongoUS.InsertUser(storage.UserFormat{
		Username: "test1",
		Password: "123",
	})
	storage.MongoUS.InsertUser(storage.UserFormat{
		Username: "test2",
		Password: "123",
	})
	storage.MongoTS.InsertToken(storage.TokenFormat{
		PasswordID: "test2",
		Token:      "token1",
	})

	testcases := []struct {
		requestPath string
		header      map[string]string
		body        string
		remoteAddr  string

		wcode            int
		wbody            string
		relatedUser      storage.UserFormat
		modifiedPassword string
	}{
		{
			"password/resetpassword",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			"username=test1&password=123&token=token1",
			"ip:port1",

			http.StatusBadRequest,
			`{"code":2000,"message":"Token已失效或不存在"}` + "\n",
			storage.UserFormat{Username: "test1"},
			"123",
		},
		{
			"password/resetpassword",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			"username=test2&password=www&token=token1",
			"ip:port1",

			http.StatusOK,
			`{"success":true}` + "\n",
			storage.UserFormat{Username: "test2"},
			"www",
		},
	}
	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	testroute := RoutesToken["ResetPassword"]
	router.Handle(testroute.Method, testroute.Pattern, testroute.HandlerFunc)

	for i, tt := range testcases {
		var url string
		url = "http://" + path.Join("exmaple.com", tt.requestPath)
		w := testutil.HandleWithRequestInfo(router, testroute.Method, url, tt.body, tt.header, tt.remoteAddr)

		if w.Code != tt.wcode {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.wcode)
			continue
		}
		if string(w.Body.Bytes()) != tt.wbody {
			t.Errorf("#%d: HTTP response body = %q, want = %q", i, string(w.Body.Bytes()), tt.wbody)
		}
		results, _ := storage.MongoUS.GetUser(tt.relatedUser, false, []string{}, []string{})
		getPassword := ""
		if len(results) != 0 {
			getPassword = results[0].Password
		}
		if getPassword != tt.modifiedPassword {
			t.Errorf("#%d: Password Mismatch! Get = %s. Want = %s", i, getPassword, tt.modifiedPassword)
		}
	}
}
