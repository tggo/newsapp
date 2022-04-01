package repository

import (
	"context"

	"boostersNews/internal/app/config"
	"boostersNews/internal/app/repository/storage/postgres"
	"boostersNews/internal/news"
	newsStorage "boostersNews/internal/news/storage/postgres"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

// Repositories provides access all Repositories store.
type Repositories struct {
	news news.Repository
}

func InitPostgres(ctx context.Context, configuration config.Databases, logger *zap.Logger) *Repositories {
	pool, err := postgres.NewConnection(ctx, configuration.PostgresURL, logger)
	if err != nil {
		panic(err)
	}

	if configuration.MigrateEnable {
		logger.Debug("need run auto migration",
			zap.String("source", configuration.MigrateDirectory))

		m, errMigration := migrate.New(configuration.MigrateDirectory, configuration.PostgresURL)
		if errMigration != nil {
			panic(errMigration)
		}
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			panic(err)

		}
	}

	runner := postgres.NewRunner(pool, logger)
	return &Repositories{
		news: newsStorage.NewRepository(runner, logger),
	}

}

func (r *Repositories) GetNews() news.Repository {
	return r.news
}
