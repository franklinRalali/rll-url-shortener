package urlshortener

import (
	"context"

	"github.com/ralali/rll-url-shortener/internal/presentations"
)

type URLShortener interface {
	ShortenURL(ctx context.Context, req presentations.ShortenURLReq) (presentations.ShortenURLRes, error)
	GetShortURL(ctx context.Context, shortCode string) (presentations.ShortURLRes, error)
	UpdateShortURL(ctx context.Context, shortCode string, req presentations.ShortURLUpdateReq) error
	DeleteShortURLByShortCode(ctx context.Context, shortCode string) error
	GetShortURLStats(ctx context.Context, shortCode string) (presentations.StatisticsRes, error)
	AddVisitCount(ctx context.Context, shortCode string) error
}
