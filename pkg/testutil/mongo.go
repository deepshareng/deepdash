package testutil

import (
	"time"

	"gopkg.in/mgo.v2"
)

func MustNewLocalSession() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:   []string{"127.0.0.1"},
		Timeout: 20 * time.Second,
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic("Mongo Dial failed")
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}
