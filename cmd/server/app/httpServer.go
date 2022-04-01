// Package app Package http server for REST-login request
package app

import (
	"context"
	"net/http"
	"time"

	"boostersNews/internal/app/config"
	ll "boostersNews/internal/app/ginx"
	"boostersNews/internal/app/services"
	newsHandler "boostersNews/internal/news/handler/http"
	"boostersNews/pkg/helpers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

type (
	// Server is http server for handing REST-auth and WebSocket
	Server struct {
		ctx      context.Context
		services *services.Services
		config   *config.Config
		logger   *zap.Logger
		httpHost string
	}
)

// NewServer create HTTP-server instance
func NewServer(ctx context.Context, services *services.Services, config *config.Config, l *zap.Logger) *Server {
	return &Server{
		ctx:      ctx,
		services: services,
		config:   config,
		httpHost: config.HTTPHostAddr,
		logger:   l.With(zap.String("server", "http")),
	}
}

// Start HTTP server
func (as *Server) Start(ctx context.Context) error {
	server := http.Server{
		Addr:    as.httpHost,
		Handler: as.SetupRoutes(),
	}
	hf := server.ListenAndServe

	as.logger.With(zap.String("host", as.httpHost)).Info("HTTP server")
	select {
	case err := <-helpers.RegisterErrorChannel(hf):
		as.logger.With(zap.Error(err)).Error("Shutdown http server")
		return server.Shutdown(ctx)
	case <-ctx.Done():
		as.logger.Info("Shutdown http server, by context.Done")
		return server.Shutdown(ctx)
	}
}

func ginWebsocketCORS() cors.Config {
	configObj := cors.DefaultConfig()
	configObj.AllowMethods = []string{http.MethodPost, http.MethodOptions, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodHead}
	configObj.AllowCredentials = true
	configObj.AllowAllOrigins = true
	configObj.AddAllowHeaders(
		"Accept",
		"Authorization",
		"Content-Type",
		"Content-Length",
		"X-CSRF-Token",
		"Token",
		"session",
		"Origin",
		"Host",
		"Connection",
		"Accept-Encoding",
		"Accept-Language",
		"X-Locale",
		"X-Lang",
	)

	return configObj
}

// SetupRoutes setup routes for HTTP server
func (as *Server) SetupRoutes() *gin.Engine {
	if as.config.Environment != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(ll.GinZap(as.logger, time.RFC3339, map[string]bool{"/metrics": true}))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(ll.RecoveryWithZap(as.logger, true))

	router.Use(cors.New(ginWebsocketCORS()))

	// metrics and debug
	ginprometheus.NewPrometheus("gin").Use(router)  // /metrics
	pprof.Register(router, "/internal/debug/pprof") // /debug/pprof

	newsHandler.Register(router, as.services.GetNews(), as.logger)
	return router
}
