package main

import (
	"context"
	"math/rand"
	"time"

	"boostersNews/internal/app/config"
	ll "boostersNews/pkg/logger"

	"github.com/chapsuk/grace"
	"go.uber.org/zap"
)

var (
	appName     string
	release     string
	gitHash     string
	buildDate   string
	buildNumber string
)

func main() {
	rand.Seed(time.Now().UnixNano()) // some time we generate random data
	ctx := grace.ShutdownContext(context.Background())

	configuration := config.NewConfig()

	logger, err := ll.New(appName, configuration.Environment, release, buildDate, buildNumber, gitHash)
	if err != nil {
		panic(err)
	}
	logger.Info("starting app",
		zap.String("gitHash", gitHash),
		zap.String("buildDate", buildDate),
		zap.String("buildNumber", buildNumber),
	)

	_ = ctx
}
