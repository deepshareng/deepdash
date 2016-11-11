package resource

import (
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/pkg/testutil"
	"github.com/MISingularity/deepdash/storage"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TestPostUserHandler(t *testing.T) {
	dbName := "test-postuserhandler"
	collName := "postuserhandler"
	tokencollName := "token"
	handler.DEBUG = true
	session := testutil.MustNewLocalSession()
	defer session.Close()
	uc := session.DB(dbName).C(collName)
	tc := session.DB(dbName).C(tokencollName)
	// defer session.DB(dbName).DropDatabase()
	storage.MongoUS, _ = storage.NewMongoUserStorageService(uc)
	storage.MongoTS, _ = storage.NewMongoTokenStorageService(tc, time.Duration(0))

	testcases := []struct {
		requestPath string
		header      map[string]string
		body        string
		remoteAddr  string

		wcode      int
		wbody      string
		insertUser storage.UserFormat
	}{
		{
			"this-is-a-clandestine-resource/users",
			map[string]string{
				"Content-Type": "application/x-www-form-urlencoded; param=value",
			},
			"username=test&password=123",
			"ip:port1",

			http.StatusOK,
			`{"value":"Already registered! Please activate your account in 15 minutes, otherwise your account will be removed!"}` + "\n",
			storage.UserFormat{
				Username: "test",
				Password: "123",
				Activate: "0",
			},
		},
	}
	var store = sessions.NewCookieStore([]byte("secret1"))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(sessions.Sessions("mysession", store))
	testroute := RoutesRegister["PostUser"]
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
		results, err := storage.MongoUS.GetUser(tt.insertUser, true, []string{}, []string{})
		if err != nil {
			t.Errorf("#%d: Get App Failed! Err Msg = %v", i, err)
			continue
		}
		if len(results) != 1 {
			t.Errorf("#%d: Get User Mismatch! Get = %d, Want = %d", i, len(results), 1)
		}
	}
}
