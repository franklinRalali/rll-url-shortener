package http

import (
	"context"

	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/internal/server"
	"github.com/ralali/rll-url-shortener/pkg/logger"
)

// Start function handler starting http listener
func Start(ctx context.Context) {

	serve := server.NewHTTPServer()
	defer serve.Done()
	logger.Info(logger.SetMessageFormat("starting rll-url-shortener services... %d", serve.Config().App.Port), logger.EventName(consts.LogEventNameServiceStarting))

	if err := serve.Run(ctx); err != nil {
		logger.Warn(logger.SetMessageFormat("service stopped, err:%s", err.Error()), logger.EventName(consts.LogEventNameServiceStarting))
	}

	return
}
