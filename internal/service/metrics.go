package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsService struct {
	httpRequestCounters *prometheus.CounterVec
}

func NewServiceMetrics() *MetricsService {
	httpRequestCounters := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "winte_http_request_counters",
		Help: "The total number of request",
	}, []string{"method", "path"})
	return &MetricsService{httpRequestCounters}
}

func (s MetricsService) CountRequest(method, path string) {
	s.httpRequestCounters.WithLabelValues(method, path).Inc()
}
