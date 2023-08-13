package urlshortener

import "context"

type URLShortener interface {
	ShortenURL(ctx context.Context, req ShortenURLReq) (ShortenURLRes, error)
	GetShortURL(ctx context.Context, shortCode string) (ShortURLRes, error)
	UpdateShortURL(ctx context.Context, shortCode string, req ShortURLUpdateReq) error
	DeleteShortURLByShortCode(ctx context.Context, shortCode string) error
	GetShortURLStats(ctx context.Context, shortCode string) (StatisticsRes, error)
	AddVisitCount(ctx context.Context, shortCode string) error
}
