package metrics

import (
	"fmt"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Monitor on charge to create a prometheus metrics for specific end point
type Monitor struct {
	Prefix       string
	EndpointName string
	// RequestTotal is a metric about how many times has been visited an endpoint
	RequestTotal *prometheus.CounterVec

	PercentilRespTime *prometheus.SummaryVec

	RespTime prometheus.Gauge
}

// NewMonitor generates a pointer to Monitor
func NewMonitor(e *gin.Engine, prefix string) *Monitor {

	p := ginprometheus.NewPrometheus(prefix)
	p.Use(e)

	m := &Monitor{Prefix: prefix}
	m.RegisterMetrics()
	return m
}

// RegisterMetrics creates the custom metrics
func (m *Monitor) RegisterMetrics() {

	m.RequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_", m.Prefix, "requests_total"),
			Help: "Total HTTP requests processed for an endpoint ",
		}, []string{"handler", "method"},
	)
	m.RespTime = promauto.NewGauge(prometheus.GaugeOpts{
		Name: fmt.Sprintf("%s_%s", m.Prefix, "response_time"),
		Help: "It generates the duration of time response for an endpoint",
	})

	m.PercentilRespTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       fmt.Sprintf("%s_%s", m.Prefix, "percentile_response_time"),
		Help:       "It generates the percentil time response for an endpoint",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"handler", "method"})

	prometheus.Register(m.RequestTotal)
	prometheus.Register(m.RespTime)
	prometheus.Register(m.PercentilRespTime)
}
