package storage

import "gopkg.in/mgo.v2"

var MongoSession *mgo.Session

var MongoUS UserStorageService

var MongoSS SequenceService

var MongoAS AppStorageService

var RemoteAppInfoAS AppStorageService

var MongoCS ChannelStorageService

var MongoASS AppSelectedItemsService

var MongoTS TokenStorageService

var Counter *mgo.Collection

var MongoSmsService SmsService
