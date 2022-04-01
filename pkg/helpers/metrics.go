package helpers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// WebsocketRequestCount total websocket request
	// sum(rate(requests_total[1m])) — общий rps
	// sum(rate(requests_total[1m])) by (handler) — rps по хендлерам
	// rate(requests_total{handler="/some/url"}[1m]) — rps одного хендлера
	// rate(requests_total{code!=200}[1m]) — rps только ошибок

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

	WebsocketRequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "websocket_request_total",
		Help: "The total number opened websocket connections",
	},
		[]string{"domain", "handler"},
	)

	GameStartedCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "poker_game_started_total",
	},
		[]string{"table"},
	)

	GameResultCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "poker_game_result_total",
	},
		[]string{"table", "who", "bet"},
	)

	WebsocketRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "websocket_request_duration_seconds",
	},
		[]string{"domain", "handler"},
	)

	PlayersAtTableGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "poker_player_at_table",
		Help: "The player per table",
	},
		[]string{"table"},
	)

	CardProcessDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "poker_card_process_duration_ms",
		Buckets: []float64{1, 5, 10, 25, 50, 75, 100, 125, 150, 200, 400, 600, 800, 1000},
	},
		[]string{"table"},
	)

	PostgresRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "pgx_query_duration_seconds",
	}, []string{"type"},
	)

	UserPingTimeDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "poker_user_ping_timeout_duration_seconds",
		Help:    "ping time of user message received",
		Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 14, 15},
	},
		[]string{"user_id"},
	)
)
