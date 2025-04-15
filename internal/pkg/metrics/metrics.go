package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// メトリクス定義.
var (
	ViewerCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "dkw",
			Name:      "viewer_count",
			Help:      "Number of viewer_count",
		}, []string{"trackName"})
)

func init() {
	registerMetrics()
}

func registerMetrics() {
	prometheus.MustRegister(ViewerCount)
}
