package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func Counter(c echo.Context) error {
	myCounter := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "my_handler_executions",
		Help:        "Counts executions of my handler function.",
		ConstLabels: prometheus.Labels{"version": "1234"},
	})
	if err := prometheus.Register(myCounter); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	myCounter.Inc()
	c.JSON(http.StatusOK, "OK")
	return nil
}

func TotalRequests(c echo.Context) error {

	var totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)
	totalRequests.WithLabelValues("200", "GET").Inc()
	totalRequests.WithLabelValues("500", "GET").Inc()
	totalRequests.WithLabelValues("503", "GET").Inc()
	if err := prometheus.Register(totalRequests); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, "OK")
	return nil
}
