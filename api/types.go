package api

import "github.com/MISingularity/deepdash/storage"

type SimpleSuccessResponse struct {
	Success bool `json:"success"`
}

type SimpleResponseData struct {
	Value string `json:"value"`
}

type ResponseData struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type PrivateUserStatus struct {
	storage.UserFormat
	AccountStatus         string `json:"accountStatus"`
	LastWeekAccountStatus string `json:"lastWeekAccountStatus"`
}

type StatusList struct {
	UserData      []PrivateUserStatus `json:"userdata"`
	AppData       []PrivateAppStatus  `json:"appdata"`
	LastWeekCount StatusCounter       `json:"lastweekcount"`
	AccountCount  StatusCounter       `json:"accountcount"`
}

type StatusCounter struct {
	Freeze          int `json:"freeze"`
	Integrating     int `json:"integrating"`
	Integrated      int `json:"integrated"`
	IntegrateFailed int `json:"integratefailed"`
	Registered      int `json:"registered"`
}

type PrivateTotalStatus struct {
	LinkDemontration int
	LinkShare        int
	AppInstall       int
	AppOpen          int
	RegisterUser     int
	AppIntegration   int
	Device           int64
}

type PrivateAppStatus struct {
	UserName         string
	SubjectApp       string
	AppName          string
	AppID            string
	LinkDemontration int
	LinkShare        int
	AppInstall       int
	AppOpen          int
	AccountStatus    string
}

type ChannelGetReturnType struct {
	Data []map[string]string `json:"data"`
}

type Smses struct {
	Smses []storage.AppSms `json:"smses"`
}

type ChannelStatisticsGetReturnType struct {
	Data []map[string]interface{} `json:"data"`
}

type AppChannelInfo struct {
	Data []storage.AppChannel `json:"data"`
}

type AppListInfo struct {
	AppID       string            `json:"appid"`
	AppName     string            `json:"appname"`
	IconUrl     string            `json:"iconurl"`
	ChannelInfo []storage.AppItem `json:"channelinfo"`
}

type AppsGetReturnType struct {
	Applist []AppListInfo `json:"applist"`
}

type TypesGetRetrunType struct {
	Typelist []Typeitem `json:"typelist"`
}

type Typeitem struct {
	TypeID   string `json:"typeid"`
	TypeName string `json:"typename"`
}

type Event struct {
	Event   string `json:"event"`
	Display string `json:"display"`
}

type Eventlist struct {
	Eventlist []Event `json:"eventlist"`
}

type AuthSuccessValues struct {
	Token       string `json:"token"`
	Username    string `json:"username"`
	Displayname string `json:"displayname"`
}

type OauthUserData struct {
	Username  string `json:"username"`
	AccountID string `json:"accountid"`
}
