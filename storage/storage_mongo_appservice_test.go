package storage

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MISingularity/deepdash/pkg/testutil"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoAppInsert(t *testing.T) {
	dbName := "testuserupdate"
	collName := "userupdate"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	as := NewMongoAppStorageService(c)

	testcases := []struct {
		insert AppFormat
	}{
		{AppFormat{
			AppID:     "appid1",
			AccountID: "account1",
			AppName:   "app1",
			PkgName:   "pkg1",
		}},
		{AppFormat{
			AppID:     "appid2",
			AccountID: "account1",
			AppName:   "app1",
			PkgName:   "pkg1",
		}},
	}
	for _, v := range testcases {
		err := as.InsertApp(v.insert)
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

func TestMongoAppGet(t *testing.T) {
	dbName := "testappget"
	collName := "appget"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	as := NewMongoAppStorageService(c)
	PrepareAppItems(as)

	testcases := []struct {
		search AppFormat
		match  AppFormat
	}{
		// get by one flied match
		{
			AppFormat{AppID: "appid1"},
			AppFormat{
				AppID:     "appid1",
				AccountID: "account1",
				AppName:   "app1",
				PkgName:   "pkg1",
			},
		},
		// get by two fileds match
		{
			AppFormat{AccountID: "account1", AppName: "app2"},
			AppFormat{
				AppID:     "appid2",
				AccountID: "account1",
				AppName:   "app2",
				PkgName:   "pkg2",
			},
		},
		// get no item match
		{
			AppFormat{AccountID: "account2", AppName: "app12"},
			AppFormat{},
		},
	}
	for _, v := range testcases {
		result, err := as.GetApp(v.search)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get app info failed, error Msg=%v\n", err)
			continue
		}
		if !reflect.DeepEqual(result, v.match) {
			t.Errorf("Update operation in App storage failed, want = %v, get = %v", v.match, result)
			continue
		}
	}
}

func TestMongoAppListGet(t *testing.T) {
	dbName := "testapplistget"
	collName := "applistget"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	as := NewMongoAppStorageService(c)
	PrepareAppItems(as)

	testcases := []struct {
		search  AppFormat
		applist map[string]AppFormat
	}{
		{
			AppFormat{AccountID: "account1"},
			map[string]AppFormat{
				"appid1": AppFormat{
					AppID:     "appid1",
					AccountID: "account1",
					AppName:   "app1",
					PkgName:   "pkg1",
				},
				"appid2": AppFormat{
					AppID:     "appid2",
					AccountID: "account1",
					AppName:   "app2",
					PkgName:   "pkg2",
				},
			},
		},
		{
			AppFormat{AccountID: "account8"},
			map[string]AppFormat{},
		},
	}
	for _, v := range testcases {
		results, _ := as.GetApplist(v.search)
		for _, result := range results {
			if !reflect.DeepEqual(v.applist[result.AppID], result) {
				t.Errorf("Get applist mismatch, item get = %v, want = %v", result, v.applist[result.AppID])
				break
			}
		}
	}
}

func TestMongoAppUpdate(t *testing.T) {
	dbName := "testappupdate"
	collName := "appupdate"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	as := NewMongoAppStorageService(c)
	PrepareAppItems(as)

	testcases := []struct {
		updateSearch AppFormat
		update       AppFormat
		upsert       bool
		match        AppFormat
	}{
		// update one filed
		{
			AppFormat{AppID: "appid1"},
			AppFormat{
				AppID:     "appid1",
				AccountID: "account1",
				AppName:   "changeappname",
				PkgName:   "pkg1",
			},
			false,
			AppFormat{
				AppID:     "appid1",
				AccountID: "account1",
				AppName:   "changeappname",
				PkgName:   "pkg1",
			},
		},
		// update one filed with upsert option
		{
			AppFormat{AppID: "appid2", AppName: "app2"},
			AppFormat{
				AppID:     "appid2",
				AccountID: "account1",
				AppName:   "app2",
				PkgName:   "changepkg",
			},
			true,
			AppFormat{
				AppID:     "appid2",
				AccountID: "account1",
				AppName:   "app2",
				PkgName:   "changepkg",
			},
		},
		// update two fields
		{
			AppFormat{AppID: "appid3"},
			AppFormat{
				AppID:     "appid3",
				AccountID: "account3",
				AppName:   "changeappname",
				PkgName:   "changepkg",
			},
			false,
			AppFormat{
				AppID:     "appid3",
				AccountID: "account3",
				AppName:   "changeappname",
				PkgName:   "changepkg",
			},
		},
		// udpate nonexist item
		{
			AppFormat{AppID: "appid13"},
			AppFormat{AppName: "changeappname"},
			false,
			AppFormat{},
		},
		// upsert
		{
			AppFormat{AppID: "appid9"},
			AppFormat{
				AppID:     "appid9",
				AccountID: "account9",
				AppName:   "app9",
				PkgName:   "pkg9",
			},
			true,
			AppFormat{
				AppID:     "appid9",
				AccountID: "account9",
				AppName:   "app9",
				PkgName:   "pkg9",
			},
		},
	}

	for _, v := range testcases {
		err := as.UpdateApp(v.updateSearch, v.update, v.upsert)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Update app failed, error Msg=%v\n", err)
			continue
		}
		result, err := as.GetApp(v.updateSearch)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get user info failed, error Msg=%v\n", err)
			continue
		}

		if !reflect.DeepEqual(result, v.match) {
			t.Errorf("Update operation in User storage failed, want = %v, get = %v", v.match, result)
			continue
		}
	}
}
