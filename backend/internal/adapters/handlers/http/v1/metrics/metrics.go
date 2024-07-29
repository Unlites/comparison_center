package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestMetric = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace: "app",
		Subsystem: "http",
		Name:      "request",
		Help:      "Request status and duration",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	},
	[]string{"status"},
)

func ObserveRequest(duration time.Duration, status int) {
	requestMetric.WithLabelValues(strconv.Itoa(status)).Observe(duration.Seconds())
}
