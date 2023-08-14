package urlshortener

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/entity"
	"github.com/ralali/rll-url-shortener/internal/presentations"
	"github.com/ralali/rll-url-shortener/internal/repositories"
	"github.com/ralali/rll-url-shortener/pkg/cache"
	"github.com/ralali/rll-url-shortener/pkg/rands"
)

const (
	originUrlCacheKeyF = "short-url-origin-%s"
)

var (
	originUrlCacheTimeout = 5 * time.Minute
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
	randGen  rands.Randoms
}

func NewURLShortener(
	urlsRepo repositories.URLs,
	cacher cache.Cacher,
	conf appctx.Config,
	randGen rands.Randoms,
) URLShortener {
	return &urlShortener{
		urlsRepo: urlsRepo,
		cacher:   cacher,
		conf:     conf,
		randGen:  randGen,
	}
}

func (u *urlShortener) ShortenURL(ctx context.Context, req presentations.ShortenURLReq) (presentations.ShortenURLRes, error) {
	var res presentations.ShortenURLRes

	shortCode, err := u.randGen.String(8, rands.DefaultCharSet)
	if err != nil {
		return res, err
	}

	url := entity.URL{
		UserID:    req.UserID,
		URL:       req.OriginURL,
		ShortCode: shortCode,
	}

	lastId, err := u.urlsRepo.Upsert(ctx, url)
	if err != nil {
		return res, err
	}

	err = u.cacher.Set(ctx, fmt.Sprintf(originUrlCacheKeyF, shortCode), req.OriginURL, originUrlCacheTimeout)
	if err != nil {
		return res, err
	}

	res = presentations.ShortenURLRes{
		ID:        lastId,
		ShortCode: shortCode,
	}

	return res, nil
}

func (u *urlShortener) GetShortURL(ctx context.Context, shortCode string) (presentations.ShortURLRes, error) {
	var res presentations.ShortURLRes

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

func (u *urlShortener) UpdateShortURL(ctx context.Context, shortCode string, req presentations.ShortURLUpdateReq) error {
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

func (u *urlShortener) DeleteShortURLByShortCode(ctx context.Context, shortCode string) error {
	// check if short url exist
	_, err := u.urlsRepo.FindOneOriginURLByShortCode(ctx, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorShortURLNotFound{ShortCode: shortCode}
		}

		return err
	}

	if err = u.urlsRepo.DeleteByShortCode(ctx, shortCode); err != nil {
		return err
	}

	// delete the origin url from cache
	err = u.cacher.Delete(ctx, fmt.Sprintf(originUrlCacheKeyF, shortCode))

	return err
}

func (u *urlShortener) GetShortURLStats(ctx context.Context, shortCode string) (presentations.StatisticsRes, error) {
	var res presentations.StatisticsRes
	visCount, err := u.urlsRepo.FindOneVisitCountByShortCode(ctx, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, ErrorShortURLNotFound{ShortCode: shortCode}
		}

		return res, err
	}

	res.VisitCount = visCount

	return res, nil
}

func (u *urlShortener) AddVisitCount(ctx context.Context, shortCode string) error {
	err := u.urlsRepo.AddVisitCountByShortCode(ctx, shortCode, 1)
	return err
}
