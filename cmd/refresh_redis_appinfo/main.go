package main

import (
	"flag"
	"os"
	"time"

	"github.com/MISingularity/deepdash/storage"
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	fs := flag.NewFlagSet("deepstatsd-reintroduce-config", flag.ExitOnError)
	appinfoAddr := fs.String("appinfo-addr", "", "Specify the appinfo address")
	mongoAddr := fs.String("mongo-addr", "", "Specify the raw data mongo database URL")
	mongoDB := fs.String("mongodb", "deepdash", "Specify the Mongo DB")
	mongocollapp := fs.String("mongocoll-app", "appcoll", "Specify the Mongo Collection for app")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if *appinfoAddr == "" || *mongoAddr == "" {
		log.Fatal("Lack appinfoAddr/mongoAddr")
	}
	session, err := mgo.DialWithTimeout(*mongoAddr, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal("Can not connect MongoDB, err:", err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	appcoll := session.DB(*mongoDB).C(*mongocollapp)
	remoteAppInfoService := storage.NewRemoteAppInfoService(*appinfoAddr)
	result := storage.AppFormat{}
	iter := appcoll.Find(bson.M{}).Iter()

	for iter.Next(&result) {
		err := remoteAppInfoService.InsertApp(result)
		if err != nil {
			log.Errorf("Remote appinfo insert failed! Err Msg=%v", err)
		}
	}
	log.Info("Finish refreshing")
	select {}
}
