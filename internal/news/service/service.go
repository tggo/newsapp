package service

import (
	"context"

	"boostersNews/internal/news"
	"boostersNews/internal/news/model"

	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, post *model.Post) (int64, error)
	Update(ctx context.Context, id int64, post *model.Post) error
	List(ctx context.Context, filter *model.Filter) ([]*model.Post, error)
	Get(ctx context.Context, id int64) (*model.Post, error)
	Delete(ctx context.Context, id int64) error
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

func (s *service) Get(ctx context.Context, id int64) (*model.Post, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.LazyDelete(ctx, id)
}

func (s *service) Create(ctx context.Context, post *model.Post) (int64, error) {
	return s.repo.Create(ctx, post)
}

func (s *service) Update(ctx context.Context, id int64, post *model.Post) error {
	return s.repo.Update(ctx, id, post)
}

func (s *service) List(ctx context.Context, filter *model.Filter) ([]*model.Post, error) {
	return s.repo.Find(ctx, filter)
}
