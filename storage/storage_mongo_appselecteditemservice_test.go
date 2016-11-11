package storage

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MISingularity/deepdash/pkg/testutil"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoAppSelectedItemInsert(t *testing.T) {
	dbName := "test-app-selecteditems-insert"
	collName := "selecteditemsinsert"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	as := NewMongoAppSelectedItemsService(c)

	testcases := []struct {
		insert AppSelectedItems
	}{
		{AppSelectedItems{
			AppID:     "appid1",
			AccountID: "acc1",
			Events:    []string{"a1", "b1", "c1"},
			Displays:  []string{"da1", "db1", "dc1"},
		}},
		{AppSelectedItems{
			AppID:     "appid2",
			AccountID: "acc2",
			Events:    []string{"a2", "b2", "c2"},
			Displays:  []string{"da2", "db2", "dc2"},
		}},
	}
	for _, v := range testcases {
		err := as.InsertAppSelectedItems(v.insert)
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

func TestMongoAppSelectedItemUpdate(t *testing.T) {
	dbName := "test-app-selecteditems-update"
	collName := "selecteditemsupdate"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	as := NewMongoAppSelectedItemsService(c)
	PrepareAppSelectedItems(as)

	testcases := []struct {
		search AppSelectedItems
		update AppSelectedItems
	}{
		{
			AppSelectedItems{
				AppID:     "appid1",
				AccountID: "acc1",
			},
			AppSelectedItems{
				AppID:     "appid1",
				AccountID: "acc1",
				Events:    []string{"a1", "b1", "c1"},
				Displays:  []string{"da1", "db1", "dc1"},
			},
		},
		{
			AppSelectedItems{
				AppID:     "appid5",
				AccountID: "acc2",
			},
			AppSelectedItems{
				AppID:     "appid5",
				AccountID: "acc2",
				Events:    []string{"a3", "b1", "c1"},
				Displays:  []string{"da3", "db1", "dc1"},
			},
		},
		{
			AppSelectedItems{
				AppID:     "appid6",
				AccountID: "acc2",
			},
			AppSelectedItems{
				AppID:     "appid6",
				AccountID: "acc2",
				Events:    []string{"a3", "b1", "c1"},
				Displays:  []string{"da3", "db1", "dc1"},
			},
		},
		{
			AppSelectedItems{
				AppID:     "appid2",
				AccountID: "acc2",
			},
			AppSelectedItems{
				AppID:     "appid2",
				AccountID: "acc2",
				Events:    []string{},
				Displays:  []string{},
			},
		},
	}

	for _, v := range testcases {
		err := as.UpdateAppSelectedItems(v.search, v.update, true)
		if err != nil {
			t.Errorf("Update App selecte failed, err msg=%v\n", err)
		}
	}
	for _, v := range testcases {
		result, err := as.GetAppSelectedItems(v.search)

		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get app selected items failed, error Msg=%v\n", err)
			continue
		}
		if !reflect.DeepEqual(result, v.update) {
			t.Errorf("Update SelectedItems operation in App storage failed, want = %v, get = %v", v.update, result)
			continue
		}
	}
}

func TestMongoAppSelectedItemsGet(t *testing.T) {
	dbName := "test-app-selecteditems-get"
	collName := "selecteditemsget"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	// defer session.DB(dbName).DropDatabase()

	as := NewMongoAppSelectedItemsService(c)
	PrepareAppSelectedItems(as)

	testcases := []struct {
		search AppSelectedItems
		match  AppSelectedItems
	}{
		// get by one flied match
		{
			AppSelectedItems{
				AppID:     "appid1",
				AccountID: "acc1",
			},
			AppSelectedItems{
				AppID:     "appid1",
				AccountID: "acc1",
				Events:    []string{"a1", "b1", "c1"},
				Displays:  []string{"da1", "db1", "dc1"},
			},
		},
		// get no item match
		{
			AppSelectedItems{
				AppID:     "appid2",
				AccountID: "acc1",
			},
			AppSelectedItems{},
		},
	}
	for _, v := range testcases {
		result, err := as.GetAppSelectedItems(v.search)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get app selected items failed, error Msg=%v\n", err)
			continue
		}
		if !reflect.DeepEqual(result, v.match) {
			t.Errorf("GetAppSelectedItems operation in App storage failed, want = %v, get = %v", v.match, result)
			continue
		}
	}
}
