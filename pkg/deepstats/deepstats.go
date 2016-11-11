package deepstats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/MISingularity/deepdash/pkg/log"
	"github.com/MISingularity/deepdash/storage"
	"github.com/MISingularity/deepshare2/deepstats/aggregate"
	"github.com/MISingularity/deepshare2/deepstats/appchannel"
	"github.com/MISingularity/deepshare2/deepstats/appevent"
	"github.com/gin-gonic/gin"
)

var (
	DEEPSTATSD_URL = "http://127.0.0.1:16759"
)

func GetLimit(c *gin.Context) (int, error) {
	var limit int
	gran := time.Hour * 24
	if c.Query("gran") != "" {
		switch c.Query("gran") {
		case "h":
			gran = time.Hour
		case "d":
			gran = time.Hour * 24
		case "w":
			gran = time.Hour * 24 * 7
		case "m":
			gran = time.Hour * 24 * 30
		case "t":
			gran = time.Hour
		}
	}

	if c.Query("start") != "" && c.Query("end") != "" {
		start, err1 := time.Parse(time.UnixDate, c.Query("start"))
		end, err2 := time.Parse(time.UnixDate, c.Query("end"))

		if err1 != nil || err2 != nil {
			log.Error("[Channel Resource] Convert start/end attribute failed!")
			return 0, fmt.Errorf("Start/End attribute format is not right.")
		}
		limit = int(end.Sub(start)/gran) + 1
	} else if c.Query("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			log.Errorf("[Channel Resource] Limit attribute format is not right.")
			return 0, fmt.Errorf("Limit attribute format is not right.")
		}
	} else {
		log.Errorf("[Channel Resource] Limit need to be direct/indirect specified.")
		return 0, fmt.Errorf("Limit need to be direct/indirect specified.")
	}
	return limit, nil
}

func GetAppEvents(appID string) (appevent.AppEvents, error) {
	log.Infof("Request GetAppEvents, URL=%s", DEEPSTATSD_URL+"/v2/appevents/"+appID+"/events")
	res, err := http.Get(DEEPSTATSD_URL + "/v2/appevents/" + appID + "/events")
	if err != nil {
		return appevent.AppEvents{}, err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return appevent.AppEvents{}, err
	}
	events := appevent.AppEvents{}
	err = json.Unmarshal(result, &events)
	if err != nil {
		return appevent.AppEvents{}, err
	}
	log.Debugf("[--Result][deepshared][GetAppEvents] Appevent=%v", events)
	return events, nil
}

func GetAppChannels(appID string) ([]storage.AppItem, error) {
	log.Infof("Request GetAppChannels, URL=%s", DEEPSTATSD_URL+"/v2/appchannels/"+appID+"/channels")
	res, err := http.Get(DEEPSTATSD_URL + "/v2/appchannels/" + appID + "/channels")
	if err != nil {
		return []storage.AppItem{}, err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return []storage.AppItem{}, err
	}
	channels := appchannel.AppChannels{}
	err = json.Unmarshal(result, &channels)
	if err != nil {
		return []storage.AppItem{}, err
	}

	channelsReturn := make([]storage.AppItem, len(channels.Channels))
	for i, _ := range channels.Channels {
		channelsReturn[i].Channelname = channels.Channels[i]
		channelsReturn[i].Typename = ""
	}

	log.Debugf("[--Result][deepshared][GetAppChannels] Appchannels=%v", channelsReturn)
	return channelsReturn, err
}

func DeleteAppChannel(appID, channel string) error {
	deepstats_url := fmt.Sprintf("%s/v2/appchannels/%s/channels?channel=%s", DEEPSTATSD_URL, appID, url.QueryEscape(channel))
	log.Infof("Request DeleteAppChannel, URL=%s", deepstats_url)
	req, err := http.NewRequest("DELETE", deepstats_url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request deepstats failed!%v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Status is unacceptable! Res status code = %d", res.StatusCode)
	}
	return nil
}

func eventSerialize(events []string, appid string, gran string, start string, end string, limit string, os string) string {
	serialize := "?appid=" + url.QueryEscape(appid) + "&gran=" + gran
	if limit != "" {

		serialize += "&limit=" + limit
	}
	if start != "" {
		serialize += "&start=" + start
	}
	if end != "" {
		serialize += "&end=" + end
	}
	if os != "" {
		serialize += "&os=" + os
	}
	for _, v := range events {
		serialize += "&event=" + v
	}
	return serialize
}

func GetChannelCounters(appid string, channel_id string, event []string, gran string, start, end string, limit string, os string) ([]*aggregate.AggregateResult, error) {
	log.Debugf("Request GetChannelCounters, URL=%s", (DEEPSTATSD_URL + "/v2/channels/" + channel_id + "/counters" + eventSerialize(event, appid, gran, start, end, limit, os)))
	res, err := http.Get(DEEPSTATSD_URL + "/v2/channels/" + channel_id + "/counters" + eventSerialize(event, appid, gran, start, end, limit, os))
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return []*aggregate.AggregateResult{}, err
	}

	channelInfo := struct {
		Counters []*aggregate.AggregateResult
	}{}

	err = json.Unmarshal(result, &channelInfo)
	if err != nil {
		return []*aggregate.AggregateResult{}, err
	}
	log.Debugf("[--Result][deepshared][GetChannelCounters] ChannelInfos=%v", aggregate.AggregateResultsToString(channelInfo.Counters))
	return channelInfo.Counters, nil
}

func GetDeviceCount(os string) (int64, error) {
	reqURL := DEEPSTATSD_URL + "/v2/device-stat/" + os
	log.Infof("Request GetChannelCounters, URL=%s", reqURL)
	res, err := http.Get(reqURL)
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return -1, err
	}

	count := map[string]int64{}
	err = json.Unmarshal(result, &count)
	if err != nil {
		return -1, err
	}

	log.Debugf("[--Result][deepshared][GetDeviceCount] DeviceCount=%v", count["Value"])
	return count["Value"], nil
}
