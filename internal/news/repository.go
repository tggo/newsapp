package news

import "context"

type Repository interface {
	Find(ctx context.Context)
}
