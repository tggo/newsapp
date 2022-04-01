package news

import (
	"context"

	"boostersNews/internal/news/model"
)

type Repository interface {
	Create(ctx context.Context, post *model.Post) (int64, error)
	Update(ctx context.Context, id int64, post *model.Post) error
	LazyDelete(ctx context.Context, id int64) error
	Find(ctx context.Context, f *model.Filter) ([]*model.Post, error)
	Get(ctx context.Context, id int64) (*model.Post, error)
}
