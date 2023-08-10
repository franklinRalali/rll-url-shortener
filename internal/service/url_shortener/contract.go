package urlshortener

import "context"

type URLShortener interface {
	ShortenURL(ctx context.Context, req ShortenURLReq) (ShortenURLRes, error)
}