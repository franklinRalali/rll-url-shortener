package repositories

import (
	"context"

	"github.com/ralali/rll-url-shortener/internal/entity"
)

const (
	TABLE_NAME_URLS = "urls"
)

type URLs interface {
	Upsert(ctx context.Context, url entity.URL) (uint64, error)
	FindOneOriginURLByShortCode(ctx context.Context, shortCode string) (string, error)
}
