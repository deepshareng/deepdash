package statistics

import (
	"github.com/MISingularity/deepdash/pkg/deepstats"
	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/MISingularity/deepdash/pkg/transmit"
	"github.com/MISingularity/deepdash/storage"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

func TotalPvCountPeriodical() {
	log.Infof("[Periodical Task - Statistics] start caculate total pv")

	periodSpan := time.Second * 20
	timer := time.NewTimer(periodSpan)

	for {
		select {
		case <-timer.C:
			log.Infof("[Periodical Task - Statistics] start caculate another round")
			totalPvCount()
			timer.Reset(periodSpan)
		}
	}
}

func QueryTotalPvCount() int {
	rst := map[string]interface{}{}
	storage.Counter.Find(bson.M{"key": "pv_total"}).One(&rst)
	return rst["pv_total"].(int)
}

func totalPvCount() {
	log.Infof("[Statistics] start caculate")
	// get app status
	appResults, err := storage.MongoAS.GetApplist(storage.AppFormat{AccountID: ""})
	if err != nil {
		return
	}

	rawEvents := []string{"url-tem-overall-value"}
	convertedEvents := transmit.RequisiteAttrs(rawEvents, "t")

	var count int = 0
	for _, v := range appResults {
		if v.AppID == "1652E90881C1FAE8" {
			continue
		}
		aggregateResult, err := deepstats.GetChannelCounters(v.AppID, "all", convertedEvents, "t", "", "", "1", "")
		if err != nil {
			return
		}
		eventRes := transmit.Calculate(rawEvents, convertedEvents, aggregateResult, "t", 1)
		count, _ = eventAddition(count, eventRes[0]["url-tem-overall-value"])
	}

	log.Infof("[Statistics] caculate result: %d", count)
	_, err = storage.Counter.Upsert(bson.M{"key": "pv_total"}, bson.M{"$set": bson.M{"pv_total": count}})
	if err != nil {
		log.Infof("[Statistics] write pv count to mongo failed: %v", err)
	}
}

func eventAddition(currentCount int, addition string) (int, error) {
	if addition == "" {
		return currentCount, nil
	}
	//increment, err := strconv.ParseInt(addition, 10, 64)
	increment, err := strconv.Atoi(addition)
	if err != nil {
		return 0, err
	}
	return currentCount + increment, nil
}
