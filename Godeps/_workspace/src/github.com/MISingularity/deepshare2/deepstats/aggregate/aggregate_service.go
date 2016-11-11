package aggregate

import (
	"reflect"
	"sort"
	"time"

	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const TimestampForm = "20140101 00:00:00 +0000 UTC"
const fakedAppIDForCompatibleColl = "all"

// AggregateService serves counter aggregation functionalities.
type AggregateService interface {
	// Store user data into a temporary storage
	//TODO Insert to different buffers for different apps
	Insert(appID string, aggregate CounterEvent) error

	// Aggregate returns the aggregated results on all records.
	// For example, if some counter records were updated in buffer,
	// and commited to persistent storage,
	// #Insert Counter(chan1, Install, 2006-01-02T10:00:00Z00:00, 1)
	// #Insert Counter(chan1, Install, 2006-01-02T12:00:00Z00:00, 1)
	// #Insert Counter(chan1, Install, 2006-01-01T12:00:00Z00:00, 1)
	// #Aggregate(chan1) => AggregateResult: [
	//    event: Install => [
	//      AggregateCount(2006-01-02T00:00:00Z00:00, 2)
	//      AggregateCount(2006-01-01T00:00:00Z00:00, 1)
	//    ]
	//  ]
	//
	// If eventFilters is empty, then all events will be returned.
	// Otherwise, it returns only the events in eventFilters.

	// QueryDuration query aggregate event in a period, and specified paritcular appid/channel/need event.
	QueryDuration(appid string, channel string, eventFilters []string, start time.Time, granularity time.Duration, limit int) ([]*AggregateResult, error)

	// Aggregate data from temporarry storage to persistent one.
	Aggregate(appID string) error

	// On the basis of user concrete aggregate granularity,
	// ConvertTimeToGranularity will regulate the aggregation
	// data time
	// E.g.
	// If user want to aggregate data by day,
	// For any time in a day, like
	// 		20060102T10:00:00Z00:00
	// It might aggreate into this time
	// 		20060102T00:00:00Z00:00
	ConvertTimeToGranularity(time.Time) time.Time
}

type generalAggregateService struct {
	buffer         map[string]*CounterEvent
	colls          map[string]*mgo.Collection
	collNamePrefix string
	mgoDb          *mgo.Database
}

func getEventMapKey(aggregate *CounterEvent) string {
	return aggregate.Timestamp.String() + "#" + aggregate.AppID + "#" + aggregate.Channel + "#" + aggregate.Event
}

func (m *generalAggregateService) queryDurationWithEnd(appid string, channel string, eventFilters []string, start time.Time, end time.Time) ([]*AggregateResultWithTimestamp, error) {
	results := []bson.M{}
	operations := buildMongoPipeOpsOfQueryAggregatedData(appid, channel, eventFilters, start, end)
	coll := GetColl(m.mgoDb, m.colls, m.collNamePrefix, appid)
	err := coll.Pipe(operations).All(&results)
	if err != nil {
		return nil, err
	}
	//If results not found, the data may be in the old collection which contains data of all appIDs
	if len(results) == 0 {
		coll := GetColl(m.mgoDb, m.colls, m.collNamePrefix, fakedAppIDForCompatibleColl)
		err := coll.Pipe(operations).All(&results)
		if err != nil {
			return nil, err
		}
	}
	aggrs := make([]*AggregateResultWithTimestamp, len(results))
	for i, result := range results {
		resCounts := result["counts"].([]interface{})
		counts := make([]*AggregateCountWithTimestamp, len(resCounts))
		for j, resCount := range resCounts {
			rc := resCount.(bson.M)
			ts := reflect.ValueOf(rc["timestamp"]).Interface().(time.Time)
			ts = ts.UTC()
			if err != nil {
				return nil, err
			}
			counts[j] = &AggregateCountWithTimestamp{
				Timestamp: ts,
				Count:     rc["count"].(int),
			}
		}
		sort.Sort(ByTimestamp(counts))
		aggrs[i] = &AggregateResultWithTimestamp{
			Event:  result["_id"].(string),
			Counts: counts,
		}
	}
	// We only record successful flow to reflect the actual workload

	return aggrs, err
}

func (m *generalAggregateService) QueryDuration(appid string, channel string, eventFilters []string, start time.Time, granulairty time.Duration, limit int) ([]*AggregateResult, error) {
	log.Infof("[AGGREGATE SERVICE] Request aggregate by appid=%s, channel=%s, eventFilter=%v, start=%s, gran=%s, limit=%d", appid, channel, eventFilters, start, granulairty, limit)
	starttime := time.Now()
	end := start.Add(granulairty * time.Duration(limit))
	results, err := m.queryDurationWithEnd(appid, channel, eventFilters, start, end)
	if err != nil {
		return []*AggregateResult{}, err
	}
	aggrs := make([]*AggregateResult, len(results))
	for i, result := range results {
		counts := make([]*AggregateCount, limit)
		lenCounts := len(result.Counts)
		interval := start
		cur := 0
		for i := 0; i < limit; i++ {
			interval = interval.Add(granulairty)
			counts[i] = &AggregateCount{}
			for cur < lenCounts && interval.After(result.Counts[cur].Timestamp) {
				counts[i].Count += result.Counts[cur].Count
				cur++
			}
		}

		aggrs[i] = &AggregateResult{
			Event:  result.Event,
			Counts: counts,
		}
	}
	// We only record successful flow to reflect the actual workload
	in.PromCounter.AggregateDuration(time.Since(starttime))
	log.Debugf("[--Result][AGGREGATE SERVICE][QueryDuration] Aggregate Results=%s, Err=%v", AggregateResultsToString(aggrs), err)
	return aggrs, err
}

func (m *generalAggregateService) Aggregate(appID string) error {
	log.Debug("[AGGREGATE SERVICE] Aggregate current buffer events")
	startTime := time.Now()
	for k, aggEvent := range m.buffer {
		if appID != fakedAppIDForCompatibleColl && appID != aggEvent.AppID {
			continue
		}
		coll := GetColl(m.mgoDb, m.colls, m.collNamePrefix, appID)
		_, err := coll.Upsert(
			bson.M{"appid": aggEvent.AppID, "channel": aggEvent.Channel, "event": aggEvent.Event, "timestamp": aggEvent.Timestamp},
			bson.M{
				"$inc": bson.M{"count": aggEvent.Count},
			},
		)

		if err != nil {
			return err
		}
		delete(m.buffer, k)
	}
	// We only record successful flow to reflect the actual workload
	in.PromCounter.AggregateDuration(time.Since(startTime))

	// return aggrs, nil
	return nil
}

func regexConstruct(events []string) []bson.RegEx {
	regexs := make([]bson.RegEx, len(events))
	for i, v := range events {
		regexs[i] = bson.RegEx{Pattern: `^` + v + `$`, Options: ""}
	}
	return regexs
}

func buildMongoPipeOpsOfQueryAggregatedData(appid string, channel string, eventFilters []string, start, end time.Time) []bson.M {
	// Filter channel_id first. Assuming that it could reduce a lot of results.
	operations := []bson.M{}

	if len(eventFilters) != 0 {
		operations = append(operations, bson.M{
			"$match": bson.M{"appid": appid, "channel": channel, "timestamp": bson.M{"$gte": start, "$lt": end}, "event": bson.M{"$in": regexConstruct(eventFilters)}},
		})
	} else {
		operations = append(operations, bson.M{
			"$match": bson.M{"appid": appid, "channel": channel, "timestamp": bson.M{"$gte": start, "$lt": end}},
		})
	}

	operations = append(operations, []bson.M{
		bson.M{
			"$group": bson.M{
				"_id": "$event",
				"counts": bson.M{
					"$push": bson.M{
						"timestamp": "$timestamp",
						"count":     "$count",
					},
				},
			},
		},
	}...)
	return operations
}

func GetColl(mgoDb *mgo.Database, colls map[string]*mgo.Collection, collNamePrefix, appid string) *mgo.Collection {
	if collNamePrefix == "" {
		log.Fatal("generalAggregateService collNamePrefix is empty")
		return nil
	}
	if mgoDb == nil {
		log.Fatal("generalAggregateService mgoDb is nil")
		return nil
	}
	if colls == nil {
		colls = make(map[string]*mgo.Collection)
	}
	if coll := colls[appid]; coll != nil {
		return coll
	}
	//for old version data, use the collection name without appID
	if appid == fakedAppIDForCompatibleColl {
		coll := mgoDb.C(collNamePrefix)
		colls[appid] = coll
		return coll
	}
	coll := mgoDb.C(collNamePrefix + "_" + appid)
	colls[appid] = coll
	return coll
}
