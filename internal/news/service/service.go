package service

import (
	"context"

	"boostersNews/internal/news"

	"go.uber.org/zap"
)

type Service interface {
	Get(ctx context.Context)
}

type service struct {
	repo   news.Repository
	logger *zap.Logger
}

func NewService(repo news.Repository, logger *zap.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger.With(zap.String("service", "news")),
	}
}

func (s *service) Get(ctx context.Context) {

}
