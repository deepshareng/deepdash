package testutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var AccountIDDebugSetting = func(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("Accountid", handler.AccountID)
	session.Set("Displayname", handler.Displayname)
	session.Save()
	c.Next()
}

// Shorthand for creating http request, asking handler to handle it,
// and returning a ResponseRecorder for checking results.
func HandleWithBody(hl http.Handler, method, url, body string) *httptest.ResponseRecorder {
	req := mustNewHTTPRequest(method, url, strings.NewReader(body))
	w := httptest.NewRecorder()
	hl.ServeHTTP(w, req)
	return w
}

// Shorthand for creating http request with header, asking handler to handle it,
// and returning a ResponseRecorder for checking results.
func HandleWithRequestInfo(hl http.Handler, method, url, body string, header map[string]string, remoteAddr string) *httptest.ResponseRecorder {
	req := mustNewHTTPRequest(method, url, strings.NewReader(body))
	req.RemoteAddr = remoteAddr
	for k, v := range header {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	hl.ServeHTTP(w, req)
	return w
}

func mustNewHTTPRequest(method, urlStr string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		log.Errorf("http.NewRequest failed!\nmethod: %s, url: %s, body: %s, err: %v", method, urlStr, body, err)
		panic(err)
	}
	return r
}

func MockResponse(code []int, body []string) (mockServer *httptest.Server, mockClient *http.Client, requestUrlHistory *([]string)) {
	index := 0
	var requestUrlHis []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUrlHis = append(requestUrlHis, r.URL.String())
		if index >= len(body) {
			fmt.Fprintln(w, "")
		} else {
			w.WriteHeader(code[index])
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, body[index])
			index++
		}
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}

	return server, httpClient, &requestUrlHis
}
