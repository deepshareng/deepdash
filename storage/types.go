package storage

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TokenFormat struct {
	PasswordID string    `bson:"password-id"`
	AccountID  string    `bson:"account-id"`
	Token      string    `bson:"token"`
	CreateAt   time.Time `bson:"createAt"`
}

type UserFormat struct {
	Id                bson.ObjectId `json:"id" bson:"_id"`
	Token             string        `json:"token" bson:"token"`
	Username          string        `json:"username" bson:"username"`
	Appname           string        `json:"appname" bson:"appname"`
	Githubname        string        `json:"githubname" bson:"githubname"`
	Password          string        `json:"password" bson:"password"`
	Phone             string        `json:"phone" bson:"phone"`
	Source            string        `json:"source" bson:"source"`
	Activate          string        `json:"activate" bson:"activate"`
	Permission        string        `json:"permission" bson:"permission"`
	Freeze            string        `json:"freeze" bson:"freeze"`
	ReturnVisitResult string        `json:"returnVisitResult" bson:"returnVisitResult"`
	ReturnVisitTime   string        `json:"returnVisitTime" bson:"returnVisitTime"`
	CreateAt          string        `json:"createat" bson:"createat"`
}

type AppFormat struct {
	AppID                         string `json:"appID" bson:"appid"`
	ShortID                       string `json:"shortID" bson:"shortid"`
	AccountID                     string `json:"accountID" bson:"accountid"`
	AppName                       string `json:"appName" bson:"appname"`
	PkgName                       string `json:"pkgName" bson:"fullpkgname"` // redundant
	IosBundler                    string `json:"iosBundler" bson:"iosbundler"`
	IosScheme                     string `json:"iosScheme" bson:"iosscheme"`
	IosDownloadUrl                string `json:"iosDownloadUrl" bson:"iosdownloadurl"`
	IosUniversalLinkEnable        string `json:"iosUniversalLinkEnable" bson:"iosunilink"` // redundant
	IosTeamID                     string `json:"iosTeamID" bson:"iosteamid"`
	IosYYBEnableBelow9            string `json:"iosYYBEnableBelow9" bson:"iosyybenablebelow9"`
	IosYYBEnableAbove9            string `json:"iosYYBEnableAbove9" bson:"iosyybenableabove9"`
	AndroidPkgname                string `json:"androidPkgname" bson:"androidpkgname"`
	AndroidScheme                 string `json:"androidScheme" bson:"androidscheme"`
	AndroidHost                   string `json:"androidHost" bson:"androidhost"`       // redundant
	AndroidAppLink                string `json:"androidAppLink" bson:"androidapplink"` // depend on androidsha256
	AndroidDownloadUrl            string `json:"androidDownloadUrl" bson:"androiddownloadurl"`
	AndroidIsDownloadDirectly     string `json:"androidIsDownloadDirectly" bson:"androidisdownloaddirectly"`
	AndroidYYBEnable              string `json:"androidYYBEnable" bson:"androidyybenable"`
	AndroidSHA256                 string `json:"androidSHA256" bson:"androidsha256"`
	YYBurl                        string `json:"yyburl" bson:"yyburl"`
	YYBenable                     string `json:"yybenable" bson:"yybenable"`                   // redundant
	AttributionPushUrl            string `json:"attributionPushUrl" bson:"attributionpushurl"` // have not known
	IconUrl                       string `json:"iconUrl" bson:"iconurl"`
	DownloadTitle                 string `json:"downloadTitle" bson:"download_title"`
	DownloadMsg                   string `json:"downloadMsg" bson:"download_msg"`
	Theme                         string `json:"theme" bson:"theme"`
	UserConfBgWeChatAndroidTipUrl string `json:"userConfBgWeChatAndroidTipUrl" bson:"userconfbgwechatandroidtipurl"`
	UserConfBgWeChatIosTipUrl     string `json:"userConfBgWeChatIosTipUrl" bson:"userconfbgwechatiostipurl"`
	ForceDownload                 string `json:"forceDownload" bson:"forcedownload"`
}

type AppItem struct {
	TypeID      string `json:"typeid" bson:"typeid"`
	Typename    string `json:"typename" bson:"typename"`
	Channelname string `json:"channelname" bson:"channelname"`
	Userdefine  string `json:"userdefine" bson:"userdefine"`
	Remark      string `json:"remark" bson:"remark"`
	MatchURL    string `json:"matchURL" bson:"matchURL"`
	Tags        []Tag  `json:"tags" bson:"tags"`
	AppID       string `json:"appid" bson:"appid"`
}

type Tag struct {
	TagID   bson.ObjectId `json:"tagid" bson:"_id"`
	TagName string        `json:"tagname" bson:"tagname"`
	TagType int           `json:"tagtype" bson:"tagtype"`
	AppID   string        `json:"appid" bson:"appid"`
}

type AppChannel struct {
	AppID       string `json:"appid" bson:"appid"`
	Channelname string `json:"channelname" bson:"channelname"`
	Channelurl  string `json:"channelurl" bson:"channelurl"`
	CreateAt    string `json:"createat" bson:"createat"`
}

type AppSms struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	ChannelName string        `json:"channelname" bson:"channelname"`
	AppID       string        `json:"appid" bson:"appid"`
	Url         string        `json:"url" bson:"url"`
	Content     string        `json:"content" bson:"content"`
	CreateAt    string        `json:"createat" bson:"createat"`
}

type AppSelectedItems struct {
	AppID     string   `json:"appid" bson:"appid"`
	AccountID string   `json:"accountid" bson:"accountid"`
	Events    []string `json:"events" bson:"events"`
	Displays  []string `json:"displays" bson:"displays"`
}

/**
 * Added from deepshare2
 */
type AppInfo struct {
	AppID     string
	ShortID   string
	AppName   string
	Android   AppAndroidInfo
	Ios       AppIosInfo
	YYBUrl    string
	YYBEnable bool
	IconUrl   string
	Theme     string
	UserConf  UserConfig
}

type AppAndroidInfo struct {
	Scheme             string
	Host               string
	Pkg                string
	DownloadUrl        string
	IsDownloadDirectly bool
	YYBEnable          bool
}

type AppIosInfo struct {
	Scheme              string
	BundleID            string
	TeamID              string
	DownloadUrl         string
	UniversalLinkEnable bool
	YYBEnableBelow9     bool
	YYBEnableAbove9     bool
	ForceDownload       bool
}

type UserConfig struct {
	BgWeChatAndroidTipUrl string
	BgWeChatIosTipUrl     string
}
