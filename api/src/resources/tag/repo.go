package tag

import (
	"context"
)

type Repository interface {
	GetTags(ctx context.Context, query string) ([]Tag, error)
}
