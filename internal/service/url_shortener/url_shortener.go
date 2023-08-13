package urlshortener

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/entity"
	"github.com/ralali/rll-url-shortener/internal/repositories"
	"github.com/ralali/rll-url-shortener/pkg/cache"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/rands"
)

const (
	originUrlCacheKeyF = "short-url-origin-%s"
)

var (
	originUrlCacheTimeout = 5 * time.Minute
)

var (
	errShortUrlNotFound = errors.New("short url not found")
)

type ErrorShortURLNotFound struct {
	ShortCode string
}

func (errNotFound ErrorShortURLNotFound) Error() string {
	return fmt.Sprintf("short URL with short code %s is not found", errNotFound.ShortCode)
}

type urlShortener struct {
	urlsRepo repositories.URLs
	cacher   cache.Cacher
	conf     appctx.Config
}

func NewURLShortener(urlsRepo repositories.URLs, cacher cache.Cacher, conf appctx.Config) URLShortener {
	return &urlShortener{
		urlsRepo: urlsRepo,
		cacher:   cacher,
		conf:     conf,
	}
}

func (u *urlShortener) ShortenURL(ctx context.Context, req ShortenURLReq) (ShortenURLRes, error) {
	var res ShortenURLRes
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

	url := entity.URL{
		UserID:    req.UserID,
		URL:       req.OriginURL,
		ShortCode: shortCode,
	}

	lastId, err := u.urlsRepo.Upsert(ctx, url)
	if err != nil {
		return res, nil
	}

	err = u.cacher.Set(ctx, fmt.Sprintf(originUrlCacheKeyF, shortCode), req.OriginURL, originUrlCacheTimeout)
	if err != nil {
		return res, err
	}

	res = ShortenURLRes{
		ID:        lastId,
		ShortCode: shortCode,
	}

	return res, nil
}

func (u *urlShortener) GetShortURL(ctx context.Context, shortCode string) (ShortURLRes, error) {
	var res ShortURLRes

	// find the origin url from cache first
	cVal, err := u.cacher.Get(ctx, fmt.Sprintf(originUrlCacheKeyF, shortCode))
	if err != nil {
		return res, err
	}

	// if found, return res
	if len(cVal) != 0 {
		res.OriginURL = string(cVal)
		res.ShortURL = fmt.Sprintf("%s/%s", u.conf.App.ShortURLHost, shortCode)

		return res, nil
	}

	originUrl, err := u.urlsRepo.FindOneOriginURLByShortCode(ctx, shortCode)
	if err != nil {
		return res, err
	}

	res.OriginURL = originUrl
	res.ShortURL = fmt.Sprintf("%s/%s", u.conf.App.ShortURLHost, shortCode)

	// store the origin url to cache
	err = u.cacher.Set(ctx, fmt.Sprintf(originUrlCacheKeyF, shortCode), originUrl, originUrlCacheTimeout)

	return res, err
}

func (u *urlShortener) UpdateShortURL(ctx context.Context, shortCode string, req ShortURLUpdateReq) error {
	// check if url with short code is exist
	_, err := u.urlsRepo.FindOneOriginURLByShortCode(ctx, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorShortURLNotFound{ShortCode: shortCode}
		}

		return err
	}

	url := entity.URL{
		URL: req.OriginURL,
	}

	if err := u.urlsRepo.UpdateByShortCode(ctx, shortCode, url); err != nil {
		return err
	}

	// update the short url on cache too
	err = u.cacher.Set(ctx, fmt.Sprintf(originUrlCacheKeyF, shortCode), req.OriginURL, originUrlCacheTimeout)
	if err != nil {
		return err
	}

	return nil
}
