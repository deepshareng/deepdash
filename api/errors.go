package api

import (
	"net/http"

	"github.com/MISingularity/deepdash/pkg/httputil"
)

var (
	// common request error
	ErrPathNotFound   = NewHTTPError(http.StatusNotFound, CodePathNotFound, false)
	ErrInternalServer = NewHTTPError(http.StatusInternalServerError, CodeInternalServer, false)
	ErrBadJSONBody    = NewHTTPError(http.StatusBadRequest, CodeBadJSONBody, false)
	ErrBadRequestBody = NewHTTPError(http.StatusBadRequest, CodeBadRequestBody, false)

	// Auth Error
	ErrAuthLoginFail    = NewHTTPError(http.StatusBadRequest, CodeAuthLoginFail, false)
	ErrAuthResourceFail = NewHTTPError(http.StatusBadRequest, CodeAuthResourceFail, false)

	// Mongo Error
	ErrMongoBrokenPipe    = NewHTTPError(http.StatusBadRequest, CodeMongoBrokenPipe, true)
	ErrMongoOperationFail = NewHTTPError(http.StatusBadRequest, CodeMongoOperationFail, false)

	ErrMailSystem = NewHTTPError(http.StatusBadRequest, CodeMailSystemError, false)

	ErrTokenInvaild                  = NewHTTPError(http.StatusBadRequest, CodeTokenInvalid, false)
	ErrSessionSave                   = NewHTTPError(http.StatusInternalServerError, CodeSessionSaveError, false)
	ErrDeepstatsConnectFailed        = NewHTTPError(http.StatusBadRequest, CodeDeepstatsConnectFailed, false)
	ErrDeepstatsGetLimitFailed       = NewHTTPError(http.StatusBadRequest, CodeDeepstatsGetLimitFailed, false)
	ErrDeepstatsGetCounterFailed     = NewHTTPError(http.StatusBadRequest, CodeDeepstatsGetCounterFailed, false)
	ErrDeepstatsGetAppChannelsFailed = NewHTTPError(http.StatusBadRequest, CodeDeepstatsGetAppChannelsFailed, false)

	ErrUserAlreadyActive = NewHTTPError(http.StatusBadRequest, CodeUserAlreadyActive, false)
	ErrUserAlreadyExist  = NewHTTPError(http.StatusBadRequest, CodeUserAlreadyExist, false)
	ErrCallBackFailed    = NewHTTPError(http.StatusBadRequest, CodeCallBackFailed, false)
	// user error that related to cookie
	ErrCookieNotFound = NewHTTPError(http.StatusNotFound, CodeCookieNotFound, false)
)

const (
	CodeBadJSONBody    = 1000
	CodeInternalServer = 1001
	CodePathNotFound   = 1002
	CodeBadRequestBody = 1003

	CodeAuthLoginFail    = 1100
	CodeAuthResourceFail = 1101

	CodeMongoBrokenPipe    = 1200
	CodeMongoOperationFail = 1201
	CodeMailSystemError    = 1300

	CodeTokenInvalid     = 2000
	CodeSessionSaveError = 2100

	CodeDeepstatsConnectFailed        = 2200
	CodeDeepstatsGetLimitFailed       = 2201
	CodeDeepstatsGetCounterFailed     = 2202
	CodeDeepstatsGetAppChannelsFailed = 2203

	CodeUserAlreadyExist  = 2300
	CodeUserAlreadyActive = 2301

	CodeCallBackFailed = 3000

	CodeCookieNotFound = 5000

	CodeParamNoneValid = 6000
)

var errMsg = map[int]string{
	CodeBadJSONBody:    "错误的JSON格式",
	CodeInternalServer: "内部服务器错误",
	CodePathNotFound:   "路径未找到或路径错误",
	CodeBadRequestBody: "请求Body格式错误",

	CodeAuthLoginFail:    "登陆失败",
	CodeAuthResourceFail: "资源验证失败",

	CodeMongoBrokenPipe:    "数据库连接超时",
	CodeMongoOperationFail: "数据操作失败",

	CodeMailSystemError: "邮件系统错误",

	CodeTokenInvalid: "Token已失效或不存在",

	CodeSessionSaveError: "Session存储错误",

	CodeUserAlreadyExist:  "用户已存在",
	CodeUserAlreadyActive: "用户已激活",

	CodeDeepstatsConnectFailed: "操作Deepstats失败",

	CodeCookieNotFound: "Faild to get cookieID for the provided device",
}

func NewHTTPError(statusCode, code int, fatal bool) httputil.HTTPError {
	return httputil.HTTPError{StatusCode: statusCode, Code: code, Message: errMsg[code], Fatal: fatal}
}
