// Package bootstrap
package bootstrap

import (
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/msg"
)

func RegistryMessage()  {
	err := msg.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.SetMessageFormat("file message multi language load error %s", err.Error()))
	}

}
