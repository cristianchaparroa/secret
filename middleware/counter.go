package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// CounterRequest generates a metric with prometheus each time that is
// made a requeust in each endpoint
func CounterRequest(counter *prometheus.CounterVec) gin.HandlerFunc {

	return func(c *gin.Context) {

		// before request
		path := removeParams(c)
		c.Next()
		// after request

		logrus.Debug("--> CounterRequest")

		labels := prometheus.Labels{}

		labels["handler"] = path
		labels["method"] = strings.ToLower(c.Request.Method)
		counter.With(labels).Inc()
		logrus.Debug("<-- CounterRequest")
	}
}
