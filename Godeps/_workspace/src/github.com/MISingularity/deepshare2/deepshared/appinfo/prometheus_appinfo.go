package appinfo

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type AppInfoInstrument interface {
	HTTPGetDuration(d time.Duration)    // overall time performing GET HTTP Request
	HTTPPutDuration(d time.Duration)    // overall time performing POST HTTP Request
	StorageGetDuration(d time.Duration) // time spent getting data from storage
	StoragePutDuration(d time.Duration) // time spent saving data to the storage
}

// PrometheusMatch holds metrics for all match methods.
type prometheusAppInfo struct {
	httpGetDuration    prometheus.Histogram
	httpPutDuration    prometheus.Histogram
	storageGetDuration prometheus.Histogram
	storagePutDuration prometheus.Histogram
}

var PrometheusForAppInfo = NewPrometheusForAppInfo()

func NewPrometheusForAppInfo() AppInfoInstrument {
	i := &prometheusAppInfo{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "appinfo",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP GET duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		httpPutDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "appinfo",
				Name:      "http_post_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP POST duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storageGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "appinfo",
				Name:      "storage_get_duration_milliseconds",
				Help:      "Bucketed histogram of storage get duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		storagePutDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepshare",
				Subsystem: "appinfo",
				Name:      "storage_save_duration_milliseconds",
				Help:      "Bucketed histogram of storage save duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpGetDuration)
	prometheus.MustRegister(i.httpPutDuration)
	prometheus.MustRegister(i.storageGetDuration)
	prometheus.MustRegister(i.storagePutDuration)

	return i
}

func (ps *prometheusAppInfo) HTTPGetDuration(d time.Duration) {
	ps.httpGetDuration.Observe(float64(d) / float64(time.Millisecond))
}
func (ps *prometheusAppInfo) HTTPPutDuration(d time.Duration) {
	ps.httpPutDuration.Observe(float64(d) / float64(time.Millisecond))
}

func (ps *prometheusAppInfo) StorageGetDuration(d time.Duration) {
	ps.storageGetDuration.Observe(float64(d) / float64(time.Millisecond))
}

func (ps *prometheusAppInfo) StoragePutDuration(d time.Duration) {
	ps.storagePutDuration.Observe(float64(d) / float64(time.Millisecond))
}
