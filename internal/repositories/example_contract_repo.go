// Package repositories
package repositories

import (
	"context"
	"database/sql"

	"github.com/ralali/rll-url-shortener/internal/entity"
)

const (
	TABLE_NAME_EXAMPLE = `example`
)

type WhereCondition struct{}

type Example interface {
	Find(ctx context.Context) ([]entity.Example, error)
	Upsert(ctx context.Context, p entity.Example) (uint64, error)
	Delete(ctx context.Context, id uint64) error
}

type DBTransaction interface {
	ExecTX(ctx context.Context, options *sql.TxOptions, fn func(context.Context, StoreTX) error) error
}

type StoreTX interface {
	Store(ctx context.Context, tableName string, entity interface{}) (id int64, err error)
	Update(ctx context.Context, tableName string, entity interface{}, whereConditions []WhereCondition) (affected int64, err error)
	StoreBulk(ctx context.Context, tableName string, entity interface{}) (int64, error)
	Upsert(ctx context.Context, tableName string, entity interface{}, onUpdate []WhereCondition) (affected int64, err error)
}