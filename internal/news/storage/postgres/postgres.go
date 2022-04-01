package postgres

import (
	"context"

	"boostersNews/internal/news"

	"go.uber.org/zap"

	"boostersNews/internal/app/repository/storage/postgres"
)

type newsStore struct {
	logger *zap.Logger
	runner postgres.IRunner
}

// NewRepository returns a new instance of a postgres newsStore repository.
func NewRepository(runner postgres.IRunner, logger *zap.Logger) news.Repository {
	return &newsStore{
		runner: runner,
		logger: logger.With(zap.String("repository", "news")),
	}
}

func (r *newsStore) Find(ctx context.Context) {

}
