package helpers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// WebsocketRequestCount total websocket request
	// sum(rate(requests_total[1m])) — Total rps
	// sum(rate(requests_total[1m])) by (handler) — rps by handlers
	// rate(requests_total{handler="/some/url"}[1m]) — rps onae handlers
	// rate(requests_total{code!=200}[1m]) — rps only errors

	HTTPRequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "The total number called http",
	},
		[]string{"domain", "path"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
	},
		[]string{"domain", "path"},
	)

	PostgresRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "pgx_query_duration_seconds",
	}, []string{"type"},
	)
)
