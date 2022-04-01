package postgres

import (
	"context"
	"database/sql"
	"time"

	"boostersNews/internal/news"
	"boostersNews/internal/news/model"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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

func (r *newsStore) Create(ctx context.Context, post *model.Post) (id int64, err error) {
	sql := `INSERT INTO public.jrnl_posts(created_at, title, body) 
			VALUES ($1, $2, $3) RETURNING id`
	err = r.runner.QueryRow(ctx, sql, func(row pgx.Row) error {
		err = row.Scan(&id)
		if err != nil {
			return err
		}
		return nil
	}, post.CreatedAt, post.Title, post.Body)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *newsStore) Update(ctx context.Context, id int64, post *model.Post) error {
	query := "UPDATE public.jrnl_posts SET updated_at = $2, title=$3, body=$4 WHERE id = $1"

	err := r.runner.Exec(ctx, query, func(ct pgconn.CommandTag) error {
		return nil
	}, id, time.Now().UTC(), post.Title, post.Body)
	if err != nil {
		if err == postgres.ErrNothingUpdate {
			return nil
		}
		r.logger.Error("could not exec sql", zap.Error(err))
		return err
	}
	return nil
}

func (r *newsStore) RemoveLazyDeleted(ctx context.Context) error {
	query := "DELETE from public.jrnl_posts  WHERE deleted_at is not null"
	err := r.runner.Exec(ctx, query, func(ct pgconn.CommandTag) error {
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *newsStore) LazyDelete(ctx context.Context, id int64) error {
	query := "UPDATE public.jrnl_posts SET deleted_at = $2 WHERE id = $1"

	err := r.runner.Exec(ctx, query, func(ct pgconn.CommandTag) error {
		return nil
	}, id, time.Now().UTC())
	if err != nil {
		if err == postgres.ErrNothingUpdate {
			return nil
		}
		r.logger.Error("could not exec sql", zap.Error(err))
		return err
	}
	return nil
}

func (r *newsStore) Find(ctx context.Context, f *model.Filter) ([]*model.Post, error) {
	query := `SELECT  id,  title, body,  created_at, updated_at
			FROM public.jrnl_posts 
			WHERE id >= $1 AND deleted_at is null
			ORDER BY id ASC
			LIMIT $2
			`

	list := make([]*model.Post, 0)
	err := r.runner.Query(ctx, query, func(row pgx.Rows) error {
		postObj := model.Post{}
		var updateAtNull sql.NullTime
		err := row.Scan(&postObj.ID, &postObj.Title, &postObj.Body, &postObj.CreatedAt, &updateAtNull)
		if err != nil {
			return err
		}

		// update at can be a null
		if updateAtNull.Valid {
			postObj.UpdatedAt = updateAtNull.Time
		}

		list = append(list, &postObj)
		return nil
	}, f.GetStart(), f.GetLimit())
	if err != nil {
		return nil, err
	}

	return list, nil

}

func (r *newsStore) Get(ctx context.Context, id int64) (*model.Post, error) {
	query := `SELECT  id,  title, body,  created_at, updated_at
			FROM public.jrnl_posts 
			WHERE id = $1 
			LIMIT 1
			`
	item := model.Post{}
	var updateAtNull sql.NullTime
	err := r.runner.QueryRow(ctx, query, func(row pgx.Row) error {
		err := row.Scan(&item.ID, &item.Title, &item.Body, &item.CreatedAt, &updateAtNull)
		if err != nil {
			return err
		}

		// update at can be a null
		if updateAtNull.Valid {
			item.UpdatedAt = updateAtNull.Time
		}

		return nil
	}, id)
	if err != nil {
		return nil, err
	}

	return &item, nil

}
