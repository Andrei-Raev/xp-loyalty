package handler

import (
	"github.com/gin-gonic/gin"
)

type MetricsService interface {
	CountRequest(method, path string)
}

type MetricsHandler struct {
	metricsService MetricsService
}

func NewMetricsHandler(metricsService MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsService: metricsService}
}

func (m MetricsHandler) WithMetrics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		path := ctx.FullPath()

		ctx.Next()

		m.metricsService.CountRequest(method, path)
	}
}
