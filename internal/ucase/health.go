package ucase

import (
	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"
)

type healthCheck struct {
}

func NewHealthCheck() contract.UseCase {
	return &healthCheck{}
}

func (u *healthCheck) Serve(*appctx.Data) appctx.Response {
	return appctx.Response{
		Name:    consts.ResponseSuccess,
		Message: "ok",
	}
}
