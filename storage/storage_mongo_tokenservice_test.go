package storage

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/MISingularity/deepdash/pkg/testutil"
	"gopkg.in/mgo.v2/bson"
)

func TestMongoTokenInsert(t *testing.T) {
	dbName := "testtokeninsert"
	collName := "tokeninsert"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	m, _ := NewMongoTokenStorageService(c, time.Duration(0))

	tests := []TokenFormat{
		TokenFormat{PasswordID: "PasswordID1"},
		TokenFormat{PasswordID: "PasswordID2", Token: "2222", CreateAt: time.Now()},
	}

	for i, tt := range tests {
		err := m.InsertToken(tt)
		if err != nil {
			t.Errorf("#%d: Insert failed: %v", i, err)
		}
	}

	// check: all records were inserted
	num, err := c.Find(bson.M{}).Count()
	if err != nil {
		t.Errorf("Connect mongodb collection failed, err msg=%v\n", err)
	}
	if num != len(tests) {
		t.Errorf("number of inserted records=%d, want=%d", num, len(tests))
	}
}

func TestMongoTokenDel(t *testing.T) {
	dbName := "testtokendel"
	collName := "tokendel"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	ts, _ := NewMongoTokenStorageService(c, time.Duration(0))
	PrepareTokeItems(ts)

	testcases := []struct {
		search TokenFormat
	}{
		{TokenFormat{PasswordID: "PasswordID1"}},
		{TokenFormat{PasswordID: "PasswordID2", Token: "token2"}},
	}

	for _, v := range testcases {
		err := ts.DelToken(v.search)
		if err != nil {
			t.Errorf("Del token info failed, error Msg=%v\n", err)
			continue
		}

		result, err := ts.GetToken(v.search, false)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get token info failed, error Msg=%v\n", err)
			continue
		}
		if len(result) != 0 {
			t.Errorf("number of del records=%d, want=%d", len(result), 0)
		}
	}
}

func TestMongoTokenGet(t *testing.T) {
	dbName := "testtokenget"
	collName := "tokenget"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	ts, _ := NewMongoTokenStorageService(c, time.Duration(0))
	PrepareTokeItems(ts)

	testcases := []struct {
		search TokenFormat
		match  []TokenFormat
	}{
		// get by one flied match
		{
			TokenFormat{PasswordID: "PasswordID1"},
			[]TokenFormat{
				TokenFormat{
					PasswordID: "PasswordID1",
					Token:      "token1",
					CreateAt:   time.Date(2015, time.March, 1, 1, 2, 3, 0, time.Local),
				},
			},
		},
		// get by two fileds match
		{
			TokenFormat{PasswordID: "PasswordID2", Token: "token2"},
			[]TokenFormat{
				TokenFormat{
					PasswordID: "PasswordID2",
					Token:      "token2",
					CreateAt:   time.Date(2015, time.March, 2, 1, 2, 3, 0, time.Local),
				},
			},
		},
		// get no item match
		{
			TokenFormat{PasswordID: "PasswordID3", Token: "token1"},
			[]TokenFormat{},
		},
	}

	for _, v := range testcases {
		result, err := ts.GetToken(v.search, false)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get token info failed, error Msg=%v\n", err)
			continue
		}
		if !reflect.DeepEqual(result, v.match) {
			t.Errorf("Update operation in token storage failed, want = %v, get = %v", v.match, result)
			continue
		}
	}
}

func TestMongoTokenUpdate(t *testing.T) {
	dbName := "testtserupdate"
	collName := "tserupdate"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	ts, _ := NewMongoTokenStorageService(c, time.Duration(0))
	PrepareTokeItems(ts)

	testcases := []struct {
		updateSearch TokenFormat
		update       TokenFormat
		upsert       bool
		match        TokenFormat
	}{
		// update one filed
		{
			TokenFormat{PasswordID: "PasswordID1"},
			TokenFormat{Token: "token2"},
			false,
			TokenFormat{
				PasswordID: "PasswordID1",
				Token:      "token2",
				CreateAt:   time.Date(2015, time.March, 1, 1, 2, 3, 0, time.Local),
			},
		},
		// update one filed with upsert option
		{
			TokenFormat{PasswordID: "PasswordID2"},
			TokenFormat{Token: "token2", CreateAt: time.Date(2016, time.March, 1, 1, 2, 3, 0, time.Local)},
			true,
			TokenFormat{
				PasswordID: "PasswordID2",
				Token:      "token2",
				CreateAt:   time.Date(2016, time.March, 1, 1, 2, 3, 0, time.Local),
			},
		},
		// udpate nonexist item
		{
			TokenFormat{PasswordID: "PasswordID8"},
			TokenFormat{Token: "token2"},
			false,
			TokenFormat{},
		},
		// upsert
		{
			TokenFormat{PasswordID: "PasswordID6"},
			TokenFormat{
				Token:    "token2",
				CreateAt: time.Date(2016, time.March, 1, 1, 2, 3, 0, time.Local),
			},
			true,
			TokenFormat{
				PasswordID: "PasswordID6",
				Token:      "token2",
				CreateAt:   time.Date(2016, time.March, 1, 1, 2, 3, 0, time.Local),
			},
		},
	}

	for _, v := range testcases {
		err := ts.UpdateToken(v.updateSearch, v.update, v.upsert)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Update token failed, error Msg=%v\n", err)
			continue
		}
		results, err := ts.GetToken(v.updateSearch, false)
		result := TokenFormat{}
		if len(results) > 0 {
			result = results[0]
		}
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get token info failed, error Msg=%v\n", err)
			continue
		}
		if !reflect.DeepEqual(result, v.match) {
			t.Errorf("Update operation in token storage failed, want = %v, get = %v", v.match, result)
			continue
		}
	}
}
