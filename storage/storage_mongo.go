package storage

import (
	"fmt"
	"time"

	"github.com/MISingularity/deepdash/pkg/bson_convert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoService struct {
	coll *mgo.Collection
}

func NewMongoTokenStorageService(c *mgo.Collection, expiredDuration time.Duration) (TokenStorageService, error) {
	index := mgo.Index{
		Key:         []string{"createAt"},
		ExpireAfter: expiredDuration,
	}
	c.DropIndex("createAt")
	err := c.EnsureIndex(index)
	return &MongoService{c}, err
}

func NewMongoUserStorageService(c *mgo.Collection) (UserStorageService, error) {
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	c.DropIndex("username")
	err := c.EnsureIndex(index)
	return &MongoService{c}, err
}

func NewMongoAppStorageService(c *mgo.Collection) AppStorageService {
	return &MongoService{c}
}

func NewMongoAppChannelService(c *mgo.Collection) ChannelStorageService {
	return &MongoService{c}
}

func NewMongoAppSelectedItemsService(c *mgo.Collection) AppSelectedItemsService {
	return &MongoService{c}
}

func NewMongoSequenceService(c *mgo.Collection) SequenceService {
	return &MongoService{c}
}

func NewMongoSmsService(c *mgo.Collection) SmsService {
	return &MongoService{c}
}

func (m *MongoService) GetToken(match TokenFormat, all bool) ([]TokenFormat, error) {
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	query := m.coll.Find(matchBSON)
	result := []TokenFormat{}
	if !all {
		var one TokenFormat
		err := query.One(&one)
		if err != nil {
			return []TokenFormat{}, err
		}
		result = append(result, one)
	} else {
		err := query.All(&result)
		if err != nil {
			return []TokenFormat{}, err
		}
	}
	return result, nil
}

// Insert User into mongo
func (m *MongoService) InsertToken(insert TokenFormat) error {
	insertBSON := bson_convert.Convert2BSON(insert, true, []string{}, []string{})
	return m.coll.Insert(insertBSON)
}

func (m *MongoService) UpdateToken(match TokenFormat, update TokenFormat, upsert bool) error {
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	updateBSON := bson_convert.Convert2BSON(update, true, []string{}, []string{})

	if !upsert {
		return m.coll.Update(matchBSON, bson.M{"$set": updateBSON})
	}
	_, err := m.coll.Upsert(matchBSON, bson.M{"$set": updateBSON})
	return err
}

func (m *MongoService) DelToken(del TokenFormat) error {
	delBSON := bson_convert.Convert2BSON(del, true, []string{}, []string{})
	return m.coll.Remove(delBSON)
}

func (m *MongoService) GetUser(match UserFormat, all bool, required []string, ignored []string) ([]UserFormat, error) {
	matchBSON := bson_convert.Convert2BSON(match, true, required, ignored)
	query := m.coll.Find(matchBSON)
	result := []UserFormat{}
	if !all {
		var one UserFormat
		err := query.One(&one)
		if err != nil {
			return []UserFormat{}, err
		}
		result = append(result, one)
	} else {
		err := query.All(&result)
		if err != nil {
			return []UserFormat{}, err
		}
	}
	return result, nil
}

// Insert User into mongo
func (m *MongoService) InsertUser(insert UserFormat) error {
	insertBSON := bson_convert.Convert2BSON(insert, true, []string{}, []string{})
	return m.coll.Insert(insertBSON)
}

func (m *MongoService) UpdateUser(match UserFormat, update UserFormat, upsert bool) error {
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	updateBSON := bson_convert.Convert2BSON(update, true, []string{}, []string{})

	if !upsert {
		return m.coll.Update(matchBSON, bson.M{"$set": updateBSON})
	}
	_, err := m.coll.Upsert(matchBSON, bson.M{"$set": updateBSON})
	return err
}

func (m *MongoService) GetApplist(match AppFormat) ([]AppFormat, error) {
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	result := []AppFormat{}
	err := m.coll.Find(matchBSON).All(&result)
	return result, err
}

func (m *MongoService) GetApplistBson(match bson.M) ([]AppFormat, error) {
	// FIXME: filter app that `deleted`
	match["status"] = bson.M{"$ne": 1}

	result := []AppFormat{}
	err := m.coll.Find(match).All(&result)
	return result, err
}

func (m *MongoService) GetApp(match AppFormat) (AppFormat, error) {
	result := AppFormat{}
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	err := m.coll.Find(matchBSON).One(&result)
	return result, err
}

func (m *MongoService) GetAppBson(match bson.M) (AppFormat, error) {
	result := AppFormat{}
	err := m.coll.Find(match).One(&result)
	return result, err
}

func (m *MongoService) InsertApp(insert AppFormat) error {
	insertBSON := bson_convert.Convert2BSON(insert, true, []string{}, []string{})
	match := AppFormat{}
	match.AppID = insert.AppID
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	st, err := m.coll.Upsert(matchBSON, bson.M{"$setOnInsert": insertBSON})
	if st.UpsertedId == nil {
		return fmt.Errorf("Exist inserted record!")
	}
	return err
}

func (m *MongoService) UpdateApp(match AppFormat, update AppFormat, upsert bool) error {
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	updateBSON := bson_convert.Convert2BSON(update, false, []string{}, []string{})
	if !upsert {
		return m.coll.Update(matchBSON, bson.M{"$set": updateBSON})
	}
	_, err := m.coll.Upsert(matchBSON, bson.M{"$set": updateBSON})
	return err
}

func (m *MongoService) UpdateAppBson(match bson.M, update bson.M, upsert bool) error {
	if !upsert {
		return m.coll.Update(match, bson.M{"$set": update})
	}
	_, err := m.coll.Upsert(match, bson.M{"$set": update})
	return err
}

func (m *MongoService) DeleteApp(appid string) error {
	return m.coll.Update(bson.M{"appid": appid}, bson.M{"$set": bson.M{"status": 1}})
}

func (m *MongoService) GetChannelList(match AppChannel) ([]AppChannel, error) {
	result := []AppChannel{}
	matchBSON := bson_convert.Convert2BSON(match, true, []string{"appid"}, []string{})
	err := m.coll.Find(matchBSON).All(&result)
	return result, err
}

func (m *MongoService) GetChannel(match AppChannel) (AppChannel, error) {
	result := AppChannel{}
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	err := m.coll.Find(matchBSON).One(&result)
	return result, err
}

func (m *MongoService) InsertChannel(insert AppChannel) error {
	insertBSON := bson_convert.Convert2BSON(insert, true, []string{}, []string{})
	return m.coll.Insert(insertBSON)
}

func (m *MongoService) UpdateChannel(match AppChannel, update AppChannel, upsert bool) error {
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	updateBSON := bson_convert.Convert2BSON(update, false, []string{}, []string{})
	if !upsert {
		return m.coll.Update(matchBSON, bson.M{"$set": updateBSON})
	}
	_, err := m.coll.Upsert(matchBSON, bson.M{"$set": updateBSON})
	return err
}

func (m *MongoService) RemoveChannel(remove AppChannel) error {
	removeBSON := bson_convert.Convert2BSON(remove, true, []string{}, []string{})
	return m.coll.Remove(removeBSON)
}

func (m *MongoService) GetAppSelectedItems(match AppSelectedItems) (AppSelectedItems, error) {
	result := AppSelectedItems{}
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	err := m.coll.Find(matchBSON).One(&result)
	return result, err
}

func (m *MongoService) UpdateAppSelectedItems(match AppSelectedItems, update AppSelectedItems, upsert bool) error {
	matchBSON := bson_convert.Convert2BSON(match, true, []string{}, []string{})
	updateBSON := bson_convert.Convert2BSON(update, false, []string{}, []string{})
	if !upsert {
		return m.coll.Update(matchBSON, bson.M{"$set": updateBSON})
	}
	_, err := m.coll.Upsert(matchBSON, bson.M{"$set": updateBSON})
	return err
}

func (m *MongoService) InsertAppSelectedItems(insert AppSelectedItems) error {
	insertBSON := bson_convert.Convert2BSON(insert, false, []string{}, []string{})
	return m.coll.Insert(insertBSON)
}

func (m *MongoService) GetAppleAppSiteAssociationInfo() ([]map[string]string, error) {
	result := []map[string]string{}
	/*
		err := m.coll.Find(
			bson.M{"teamid": bson.M{"$ne": ""}},
		).All(&result)
	*/
	err := m.coll.Find(bson.M{
		"$and": []bson.M{
			bson.M{"iosteamid": bson.M{"$ne": ""}},
			bson.M{"iosteamid": bson.M{"$exists": true}},
			bson.M{"iosbundler": bson.M{"$ne": ""}},
			bson.M{"iosbundler": bson.M{"$exists": true}},
		}}).Select(bson.M{
		"iosteamid": 1, "iosbundler": 1, "appid": 1, "shortid": 1,
	}).All(&result)
	return result, err
}

func (m *MongoService) GetSequenceId() (int, error) {
	// Must init sequence with:
	// db.seqcoll.insert({"seq_key": "app", "seq_val": NumberInt(1000)})
	// `1000` is the start point
	result := map[string]interface{}{}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq_val": 1}},
		ReturnNew: true,
	}
	_, err := m.coll.Find(bson.M{"seq_key": "app"}).Apply(change, &result)
	return result["seq_val"].(int), err
}

func (m *MongoService) GetSmsList(match bson.M) ([]AppSms, error) {
	result := []AppSms{}
	err := m.coll.Find(match).All(&result)
	return result, err
}

func (m *MongoService) InsertSms(insert bson.M) error {
	return m.coll.Insert(insert)
}

func (m *MongoService) RemoveSms(match bson.M) error {
	return m.coll.Remove(match)
}

func (m *MongoService) UpdateSms(appid string, smsid string, content string) error {
	return m.coll.Update(bson.M{"appid": appid, "_id": bson.ObjectIdHex(smsid)}, bson.M{"$set": bson.M{"content": content}})
}
