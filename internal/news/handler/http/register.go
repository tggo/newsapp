package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"boostersNews/internal/news/service"
)

// Register create http handler for news
func Register(router gin.IRouter, service service.Service, logger *zap.Logger) {
	h := NewHandler(service, logger)

	router.POST("posts", h.create)
	router.GET("posts", h.list)

	router.PUT("posts/:id", h.update)
	router.DELETE("posts/:id", h.delete)
}

// NewHandler create handler instance
func NewHandler(service service.Service, logger *zap.Logger) *handler {
	return &handler{
		service: service,
		logger:  logger.With(zap.String("handler", "news")),
	}
}

type handler struct {
	logger  *zap.Logger
	service service.Service
}
