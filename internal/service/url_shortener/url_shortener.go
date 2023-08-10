package urlshortener

import (
	"context"

	"github.com/ralali/rll-url-shortener/internal/repositories"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/rands"
)

type urlShortener struct {
	urlsRepo repositories.URLs
}

func NewURLShortener(urlsRepo repositories.URLs) URLShortener {
	return &urlShortener{
		urlsRepo: urlsRepo,
	}
}

func (u *urlShortener) ShortenURL(ctx context.Context, req ShortenURLReq) (ShortenURLRes, error) {
	var lf []logger.Field
	lf = append(
		lf, 
		logger.Any("user_id", req.UserID),
		logger.Any("origin_url", req.OriginURL),
	)

	shortCode, err := rands.RandString(8)
	if err != nil {
		logger.ErrorWithContext(ctx, nil, lf...)
	}

	return ShortenURLRes{}, nil
}
