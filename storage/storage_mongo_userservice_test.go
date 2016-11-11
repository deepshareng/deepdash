package storage

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MISingularity/deepdash/pkg/testutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type person struct {
	Name  string
	Phone string
}

// TestMongoSetup assumes that mongod is running in background
// and tests some basic functions of mgo library.
func TestMongoSetup(t *testing.T) {
	var err error
	testDBName := "testmongo"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	// "people" collection
	c := session.DB(testDBName).C("people")

	// Prepare phase:
	// - Insert some entries into collection
	// We can later reuse partial information of inserted entries to test query.
	persons := []*person{
		&person{"Ale", "+55 53 8116 9639"},
		&person{"Cla", "+55 53 8402 8510"},
	}
	// mgo Insert method accepts `interface{}` variadic arugments.
	personsForInsert := make([]interface{}, len(persons))
	for i, p := range persons {
		personsForInsert[i] = p
	}
	err = c.Insert(personsForInsert...)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		p   *person
		err error
	}{
		{persons[0], nil},
		{persons[1], nil},
		{&person{"No Such Person", ""}, mgo.ErrNotFound},
	}

	for i, tt := range tests {
		result := &person{}
		err = c.Find(bson.M{"name": tt.p.Name}).One(result)
		if err != nil {
			if err != tt.err {
				t.Errorf("#%d: err=%v, want=%s", i, err, tt.err)
			}
			continue
		}
		if !reflect.DeepEqual(result, tt.p) {
			t.Errorf("#%d: person=%#v, want=%#v", i, result, tt.p)
		}
	}

	err = session.DB(testDBName).DropDatabase()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMongoUserInsert(t *testing.T) {
	dbName := "testuserinsert"
	collName := "userinsert"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	m, _ := NewMongoUserStorageService(c)

	tests := []UserFormat{
		{Username: "username"},
		{Username: "username1"},
	}

	for i, tt := range tests {
		err := m.InsertUser(tt)
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

func TestMongoUserGet(t *testing.T) {
	dbName := "testuserget"
	collName := "userget"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	us, _ := NewMongoUserStorageService(c)
	PrepareUserItems(us)

	testcases := []struct {
		search UserFormat
		match  []UserFormat
	}{
		// get by one flied match
		{
			UserFormat{Username: "test1"},
			[]UserFormat{
				UserFormat{
					Username:   "test1",
					Password:   "123",
					Githubname: "git1",
					Activate:   "1",
				},
			},
		},
		// get by two fileds match
		{
			UserFormat{Username: "test2", Githubname: "git2"},
			[]UserFormat{
				UserFormat{
					Username:   "test2",
					Password:   "123",
					Githubname: "git2",
					Activate:   "1",
				},
			},
		},
		// get no item match
		{
			UserFormat{Username: "test3", Githubname: "git4"},
			[]UserFormat{},
		},
	}

	for _, v := range testcases {
		result, err := us.GetUser(v.search, false, []string{}, []string{})
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get user info failed, error Msg=%v\n", err)
			continue
		}
		for i := 0; i < len(result); i++ {
			result[i].Id = ""
		}
		if !reflect.DeepEqual(result, v.match) {
			t.Errorf("Update operation in User storage failed, want = %v, get = %v", v.match, result)
			continue
		}
	}
}

func TestMongoUserUpdate(t *testing.T) {
	dbName := "testuserupdate"
	collName := "userupdate"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	c := session.DB(dbName).C(collName)
	defer session.DB(dbName).DropDatabase()

	us, _ := NewMongoUserStorageService(c)
	PrepareUserItems(us)

	testcases := []struct {
		updateSearch UserFormat
		update       UserFormat
		upsert       bool
		match        UserFormat
	}{
		// update one filed
		{
			UserFormat{Username: "test1"},
			UserFormat{Password: "111"},
			false,
			UserFormat{
				Username:   "test1",
				Password:   "111",
				Githubname: "git1",
				Activate:   "1",
			},
		},
		// update one filed with upsert option
		{
			UserFormat{Username: "test2"},
			UserFormat{Password: "111"},
			true,
			UserFormat{
				Username:   "test2",
				Password:   "111",
				Githubname: "git2",
				Activate:   "1",
			},
		},
		// update two fields
		{
			UserFormat{Username: "test2"},
			UserFormat{Password: "111", Githubname: "gitchange"},
			false,
			UserFormat{
				Username:   "test2",
				Password:   "111",
				Githubname: "gitchange",
				Activate:   "1",
			},
		},
		// udpate nonexist item
		{
			UserFormat{Username: "test5"},
			UserFormat{Password: "123"},
			false,
			UserFormat{},
		},
		// upsert
		{
			UserFormat{Username: "test4"},
			UserFormat{Password: "123", Githubname: "git4"},
			true,
			UserFormat{
				Username:   "test4",
				Password:   "123",
				Githubname: "git4",
			},
		},
	}

	for _, v := range testcases {
		err := us.UpdateUser(v.updateSearch, v.update, v.upsert)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Update user failed, error Msg=%v\n", err)
			continue
		}
		results, err := us.GetUser(v.updateSearch, false, []string{}, []string{})
		result := UserFormat{}
		if len(results) > 0 {
			result = results[0]
		}
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Get user info failed, error Msg=%v\n", err)
			continue
		}
		result.Id = ""
		if !reflect.DeepEqual(result, v.match) {
			t.Errorf("Update operation in User storage failed, want = %v, get = %v", v.match, result)
			continue
		}
	}
}
