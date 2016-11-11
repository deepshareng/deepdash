package storage

import (
	"fmt"

	"github.com/surge/cityhash"
	"gopkg.in/mgo.v2/bson"
)

type TokenStorageService interface {
	GetToken(match TokenFormat, all bool) ([]TokenFormat, error)
	InsertToken(insert TokenFormat) error
	UpdateToken(match TokenFormat, update TokenFormat, upsert bool) error
	DelToken(del TokenFormat) error
}

type UserStorageService interface {
	GetUser(match UserFormat, all bool, required []string, ignored []string) ([]UserFormat, error)
	InsertUser(UserFormat) error
	UpdateUser(match UserFormat, update UserFormat, upsert bool) error
}

type AppStorageService interface {
	GetApplist(match AppFormat) ([]AppFormat, error)
	GetApp(match AppFormat) (AppFormat, error)
	InsertApp(AppFormat) error
	UpdateApp(match AppFormat, update AppFormat, upsert bool) error
	GetApplistBson(match bson.M) ([]AppFormat, error)
	GetAppBson(match bson.M) (AppFormat, error)
	UpdateAppBson(match bson.M, update bson.M, upsert bool) error
	GetAppleAppSiteAssociationInfo() ([]map[string]string, error)
	DeleteApp(appid string) error
}

type ChannelStorageService interface {
	GetChannelList(match AppChannel) ([]AppChannel, error)
	GetChannel(match AppChannel) (AppChannel, error)
	InsertChannel(AppChannel) error
	UpdateChannel(match AppChannel, update AppChannel, upsert bool) error
	RemoveChannel(remove AppChannel) error
}

type AppSelectedItemsService interface {
	GetAppSelectedItems(match AppSelectedItems) (AppSelectedItems, error)
	UpdateAppSelectedItems(match AppSelectedItems, update AppSelectedItems, upsert bool) error
	InsertAppSelectedItems(insert AppSelectedItems) error
}

type SequenceService interface {
	GetSequenceId() (int, error)
}

type SmsService interface {
	GetSmsList(match bson.M) ([]AppSms, error)
	InsertSms(bson.M) error
	RemoveSms(match bson.M) error
	UpdateSms(string, string, string) error
}

func GenerateID(unique string) string {
	input := []byte(unique)
	// CityHash64 return 64 bit, maxium 16bit in hex, fix width all 16bits
	return fmt.Sprintf("%016x", cityhash.CityHash64(input, uint32(len(input))))
}
