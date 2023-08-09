// Package migration
package migration

import (
	"time"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/mariadb"
)

func MigrateDatabase() {
	cfg, e := appctx.NewConfig()

	if e != nil {
		logger.Fatal(e)
	}

	mariadb.DatabaseMigration(&mariadb.Config{
		Host:         cfg.WriteDB.Host,
		Port:         cfg.WriteDB.Port,
		Name:         cfg.WriteDB.Name,
		User:         cfg.WriteDB.User,
		Password:     cfg.WriteDB.Pass,
		Charset:      cfg.WriteDB.Charset,
		Timeout:      time.Duration(cfg.WriteDB.TimeoutSecond) * time.Second,
		MaxIdleConns: cfg.WriteDB.MaxIdle,
		MaxOpenConns: cfg.WriteDB.MaxOpen,
		MaxLifetime:  time.Duration(cfg.WriteDB.MaxLifeTimeMS) * time.Millisecond,
	})
}
