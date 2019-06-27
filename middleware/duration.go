package middleware

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// PercentileResponseTime generates a metric about how long does the response
// take for a specific endpoint
func PercentileResponseTime(responseDuration *prometheus.SummaryVec) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()

		logrus.Debug("--> PercentileResponseTime")
		path := removeParams(c)
		// before request

		c.Next()
		// after request

		r := new(big.Int)
		fmt.Println(r.Binomial(1000, 10))

		elapsed := time.Since(start)
		elapsedF := float64(elapsed) / float64(time.Millisecond)

		labels := prometheus.Labels{}

		labels["handler"] = path
		labels["method"] = strings.ToLower(c.Request.Method)

		responseDuration.With(labels).Observe(elapsedF)
		logrus.Debug("<-- PercentileResponseTime")
	}
}

// ResponseTime generates a metric about how long does the response
// take for a specific endpoint
func ResponseTime(respTime prometheus.Gauge) gin.HandlerFunc {

	return func(c *gin.Context) {
		logrus.Debug("--> ResponseTime")

		timer := prometheus.NewTimer(prometheus.ObserverFunc(respTime.Set))
		defer timer.ObserveDuration()

		// before request
		c.Next()
		// after request
		logrus.Debug("<-- ResponseTime")
	}
}
