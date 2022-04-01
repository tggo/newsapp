package app

import (
	"context"

	"boostersNews/internal/app/config"

	"go.uber.org/zap"

	services "boostersNews/internal/app/services"
)

// App global app struct
type App struct {
	logger      *zap.Logger
	Config      *config.Config
	allServices *services.Services
	isError     chan struct{}
}

func New(config *config.Config, allServices *services.Services, logger *zap.Logger) *App {
	return &App{
		Config:      config,
		allServices: allServices,
		logger:      logger.With(zap.String("app", "root")),
		isError:     make(chan struct{}),
	}
}

// Run starting main app
func (a *App) Run(ctx context.Context) {
	ctx, cancelFunc := context.WithCancel(ctx)

	// // Socket server
	// socketIOServerCh := make(chan error, 1)
	// websocketServer := wsocketio.NewServer(ctx, a.allServices, a.logger)
	// go func() {
	// 	socketIOServerCh <- websocketServer.Start(ctx)
	// 	a.isError <- struct{}{}
	// }()

	// HTTP server
	httpServerErrCh := make(chan error, 1)
	httpServerObj := NewServer(ctx, a.allServices, a.Config, a.logger)
	go func() {
		httpServerErrCh <- httpServerObj.Start(ctx)
		a.isError <- struct{}{}
	}()

	<-a.isError
	cancelFunc()

	// wait graceful stopped when many services
	// select {
	// case <-socketIOServerCh:
	// 	<-httpServerErrCh
	// case <-httpServerErrCh:
	// 	<-socketIOServerCh
	// }
}
