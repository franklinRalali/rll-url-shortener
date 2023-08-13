package repositories

import (
	"context"

	"github.com/ralali/rll-url-shortener/internal/entity"
	"github.com/ralali/rll-url-shortener/pkg/mariadb"
)

type urls struct {
	db mariadb.Adapter
}

func NewUrls(db mariadb.Adapter) URLs {
	return &urls{
		db: db,
	}
}

func (r *urls) Upsert(ctx context.Context, url entity.URL) (uint64, error) {
	q := `INSERT INTO urls (user_id, origin_url, short_code) VALUES (?, ?, ?)`
	res, err := r.db.Exec(ctx, q, url.UserID, url.URL, url.ShortCode)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastId), nil
}

func (r *urls) FindOneOriginURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	q := `SELECT origin_url FROM urls WHERE short_code = ?`
	var originUrl string
	err := r.db.FetchRow(ctx, &originUrl, q, shortCode)

	return originUrl, err
}