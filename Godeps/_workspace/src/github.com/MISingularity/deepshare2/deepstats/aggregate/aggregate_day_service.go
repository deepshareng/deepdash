package aggregate

import (
	"time"

	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
)

type dayAggregateService struct {
	generalAggregateService
}

func NewDayAggregateService(db *mgo.Database, collPrefix string) AggregateService {
	return &dayAggregateService{
		generalAggregateService{
			buffer:         map[string]*CounterEvent{},
			mgoDb:          db,
			collNamePrefix: collPrefix,
			colls:          make(map[string]*mgo.Collection),
		},
	}
}

func (m *dayAggregateService) ConvertTimeToGranularity(t time.Time) time.Time {
	t = t.Local()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func (m *dayAggregateService) Insert(appID string, aggregate CounterEvent) error {
	log.Infof("[AGGREGATE SERVICE][DAY] Insert aggregate event appid: %s, event:%s", aggregate.AppID, aggregate.Event)
	log.Debugf("[AGGREGATE SERVICE][DAY] Insert aggregate event %v", aggregate)
	aggregate.Timestamp = m.ConvertTimeToGranularity(aggregate.Timestamp)
	_, ok := m.buffer[getEventMapKey(&aggregate)]
	if ok {
		m.buffer[getEventMapKey(&aggregate)].Count += aggregate.Count
	} else {
		m.buffer[getEventMapKey(&aggregate)] = &aggregate
	}
	return nil
}
