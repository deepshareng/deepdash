package instrumentation

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var PromCounter = newPromCounter()

type Counter interface {
	HTTPGETDuration(d time.Duration)
	AggregateDuration(d time.Duration)
}
type promCounter struct {
	httpGetDuration prometheus.Histogram
	agDuration      prometheus.Histogram
}

func newPromCounter() Counter {
	i := &promCounter{
		httpGetDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepstats",
				Subsystem: "counter",
				Name:      "http_get_duration_milliseconds",
				Help:      "Bucketed histogram of HTTP GET duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
		agDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "deepstats",
				Subsystem: "counter",
				Name:      "aggregate_duration_milliseconds",
				Help:      "Bucketed histogram of successful Aggregate action duration.",
				// 0.5ms -> 1000ms
				Buckets: prometheus.ExponentialBuckets(0.5, 2, 12),
			}),
	}

	prometheus.MustRegister(i.httpGetDuration)
	prometheus.MustRegister(i.agDuration)

	return i
}

func (p *promCounter) HTTPGETDuration(d time.Duration) {
	p.httpGetDuration.Observe(float64(d / time.Millisecond))
}

func (p *promCounter) AggregateDuration(d time.Duration) {
	p.agDuration.Observe(float64(d / time.Millisecond))
}
