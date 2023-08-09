// Package bootstrap
package bootstrap

import (
	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/util"
)

func RegistryLogger(cfg *appctx.Config) {

	lc := logger.Config{
		URL:         cfg.Logger.URL,
		Environment: util.EnvironmentTransform(cfg.App.Env),
		Debug:       cfg.App.Debug,
		Level:       cfg.Logger.Level,
	}

	h, e := logger.NewSentryHook(lc)

	if e != nil {
		logger.Fatal("log sentry failed to initialize Sentry")
	}

	logger.Setup(lc, h)
}
