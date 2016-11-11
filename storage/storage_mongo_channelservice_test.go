package storage

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MISingularity/deepdash/pkg/testutil"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoChannelInsert(t *testing.T) {
	dbName := "testchannel"
	collName := "channel"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	as := NewMongoAppChannelService(c)

	testcases := []struct {
		insert AppChannel
	}{
		{AppChannel{
			AppID:       "appid1",
			Channelname: "channelname1",
		}},
		{AppChannel{
			AppID:       "appid1",
			Channelname: "channelname2",
		}},
	}

	for _, v := range testcases {
		err := as.InsertChannel(v.insert)
		if err != nil {
			t.Errorf("Insert App failed, err msg=%v\n", err)
		}
	}
	num, err := c.Find(bson.M{}).Count()
	if err != nil {
		t.Errorf("Connect mongodb collection failed, err msg=%v\n", err)
	}
	if num != len(testcases) {
		t.Fatalf("number of inserted records=%d, want=%d", num, len(testcases))
	}
}

func TestMongoChannelGet(t *testing.T) {
	dbName := "testchannel"
	collName := "channelget"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	as := NewMongoAppChannelService(c)
	PrepareChannelItems(as)

	testcases := []struct {
		search AppChannel
		match  AppChannel
	}{
		// get by one flied match
		{
			AppChannel{AppID: "appid1"},
			AppChannel{
				AppID:       "appid1",
				Channelname: "channelname1",
			},
		},
		// get no item match
		{
			AppChannel{AppID: "app12"},
			AppChannel{},
		},
	}
	for _, v := range testcases {
		result, err := as.GetChannel(v.search)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get channel failed, error Msg=%v\n", err)
			continue
		}
		if !reflect.DeepEqual(result, v.match) {
			t.Errorf("GetChannel operation in App storage failed, want = %v, get = %v", v.match, result)
			continue
		}
	}
}
