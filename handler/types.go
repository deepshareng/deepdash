package handler

import "github.com/gin-gonic/gin"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

var (
	DEBUG = false
)

var AccountID string
var Displayname string

type BaseViewURL struct {
	IndexURL         string
	AddAppURL        string
	AddAppSubURL     string
	ProfileURL       string
	MarketingURL     string
	LoginURL         string
	RegisterURL      string
	ForgetURL        string
	ResetURL         string
	AppInfoURL       string
	OutstandingURL   string
	CallbackURL      string
	DeepshareURL     string
	AddAppByStepsURL string
	PrURL            string
	IntegrateURL     string
}

var BaseURL = BaseViewURL{
	"/index",
	"/addapp",
	"/addapp2",
	"/profile",
	"/marketing",
	"/login",
	"/register",
	"/forget",
	"/resetpassword",
	"/appinfo",
	"/outstanding",
	"/callback",
	"/",
	"/addapp-by-steps",
	"/pr",
	"/i",
	//"/integrate/*_", // for react browerhistory, all prefix with /integrate response share html
}

var LoginURLAuthFail string = BaseURL.LoginURL + "?error=1"
var LoginURLActivateFail = BaseURL.LoginURL + "?error=2"
var LoginURLActivateSuccess = BaseURL.LoginURL + "?success=2"
