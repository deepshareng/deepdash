package main

import (
	"flag"

	"os"
	"time"

	"github.com/MISingularity/deepdash/config"
	"github.com/MISingularity/deepdash/handler"
	"github.com/MISingularity/deepdash/handler/client"
	"github.com/MISingularity/deepdash/pkg/deepstats"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/MISingularity/deepdash/router"
	"github.com/MISingularity/deepdash/storage"
	"github.com/RangelReale/osincli"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

const (
	ENV_DEBUG     = "DEBUG"
	ENV_ACCOUNTID = "ACCOUNTID"

	PRODUCT_GITHUB_CLIENT_ID     = "9487562b91cf0e58a7f5"
	PRODUCT_GITHUB_CLIENT_SECRET = "e7de8b20bdc18a0d4c221a319ef1a585b3c187a4"

	//Development env
	DEBUG_GITHUB_CLIENT_ID     = "9b85b1e3ef194767dcde"
	DEBUG_GITHUB_CLIENT_SECRET = "f1d878b4f65c03201168e9b4c498d6f9a83e72ee"
)

var store = sessions.NewCookieStore([]byte("secret"))

func StartServer() {
	// Get necessary settings
	fs := flag.NewFlagSet("name", flag.ExitOnError)

	appinfoAddr := fs.String("appinfo-addr", "", "Specify the appinfo address")
	deepstatsAddr := fs.String("deepstats-addr", "", "Sepcify the deepstas api address")
	deepshareAddr := fs.String("deepshare-addr", "/", "Specify the deepshare main page address")
	mongoAddr := fs.String("mongo-addr", "", "Specify the raw data mongo database URL")
	mongoDBUser := fs.String("mongodb-user", "deepdash", "Specify the Mongo DB for user")
	mongocolluser := fs.String("mongocoll-user", "usercoll", "Specify the Mongo Collection for user")
	mongoDBApp := fs.String("mongodb-app", "deepdash", "Specify the Mongo DB for app")
	mongocollapp := fs.String("mongocoll-app", "appcoll", "Specify the Mongo Collection for app")
	mongoDBChannel := fs.String("mongodb-channel", "deepdash", "Specify the Mongo DB for channel")
	mongocollchannel := fs.String("mongocoll-channel", "channelcoll", "Specify the Mongo Collection for app's channel")
	mongoDBSel := fs.String("mongodb-appselection", "deepdash", "Specify the Mongo DB for app selection")
	mongocollsel := fs.String("mongocoll-sel", "selcoll", "Specify the Mongo Collection for selected items' channel")
	mongoDBToken := fs.String("mongodb-token", "deepdash", "Specify the Mongo DB for token")
	mongocolltoken := fs.String("mongocoll-token", "tokencoll", "Specify the Mongo Collection for token's channel")
	mongocollseq := fs.String("mongocoll-seq", "seqcoll", "Specify the Mongo Collection for user")
	debug := fs.Bool("debug", false, "Specify the start mode")

	// Authentication service config
	AuthAddr := fs.String("auth-addr", "http://localhost:10032", "Specify the auth address.")
	clientId := fs.String("client-id", "deepdash-local", "Specify the application client id for oauth")
	clientSecret := fs.String("client-secret", "deepdash-deepshare", "Specify the application client secret for oauth")
	redirectURL := fs.String("redirecturi", "http://localhost:10033/auth/code", "Specify the oauth redirect uri")
	// create client
	log.InitLog("[DeepDash]", "", log.LevelDebug)
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	config.Cliconfig = &osincli.ClientConfig{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		AuthorizeUrl: *AuthAddr + "/auth/authorize-token",
		TokenUrl:     *AuthAddr + "/auth/token-authentication",
		RedirectUrl:  *redirectURL,
	}
	if *mongoAddr == "" || *deepstatsAddr == "" || *appinfoAddr == "" {
		log.Info("mongo/deepshared/appinfo addr is not set!")
		os.Exit(1)
	}

	log.Infof("Set deepstas-addr as %s, deepshare-addr as %s, appinfo-addr as %s", *deepstatsAddr, *deepshareAddr, *appinfoAddr)
	deepstats.DEEPSTATSD_URL = *deepstatsAddr
	handler.BaseURL.DeepshareURL = *deepshareAddr
	handler.DEBUG = false
	// Initialize services
	session, err := mgo.DialWithTimeout(*mongoAddr, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal("Can not connect MongoDB, err:", err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	usercoll := session.DB(*mongoDBUser).C(*mongocolluser)
	appcoll := session.DB(*mongoDBApp).C(*mongocollapp)
	seqcoll := session.DB(*mongoDBApp).C(*mongocollseq)
	smscoll := session.DB(*mongoDBApp).C("sms")
	countcoll := session.DB(*mongoDBApp).C("count")

	channelcoll := session.DB(*mongoDBChannel).C(*mongocollchannel)
	selcoll := session.DB(*mongoDBSel).C(*mongocollsel)
	tokencoll := session.DB(*mongoDBToken).C(*mongocolltoken)
	log.Infof("Listening user storage %s/%s, app storage %s/%s , channel storage %s/%s", *mongoDBUser, *mongocolluser, *mongoDBApp, *mongocollapp, *mongoDBChannel, *mongocollchannel)
	storage.MongoUS, err = storage.NewMongoUserStorageService(usercoll)
	if err != nil {
		log.Fatal(err)
		return
	}
	storage.MongoSession = session.Copy()
	storage.MongoSS = storage.NewMongoSequenceService(seqcoll)
	storage.MongoSmsService = storage.NewMongoSmsService(smscoll)
	storage.MongoAS = storage.NewMongoAppStorageService(appcoll)
	storage.MongoCS = storage.NewMongoAppChannelService(channelcoll)
	storage.MongoASS = storage.NewMongoAppSelectedItemsService(selcoll)
	storage.MongoTS, err = storage.NewMongoTokenStorageService(tokencoll, time.Hour*24)
	storage.Counter = countcoll
	if err != nil {
		log.Fatal(err)
		return
	}
	storage.RemoteAppInfoAS = storage.NewRemoteAppInfoService(*appinfoAddr)
	// Set HTTP Router
	client.AuthURL = *AuthAddr

	mode := gin.ReleaseMode

	if *debug == true {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)
	ginrouter := gin.Default()
	ginrouter.Use()
	ginrouter.Use(sessions.Sessions("mysession", store))
	ginrouter.Use(gin.ErrorLogger())
	router.AddRouterByExistRule(ginrouter)
	log.Info("Listening port 10033")
	log.Fatal(ginrouter.Run(":10033"))
}

func main() {
	StartServer()
}
