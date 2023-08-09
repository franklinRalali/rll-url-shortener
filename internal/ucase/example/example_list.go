// Package example
package example

import (
	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/internal/repositories"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"

	"github.com/ralali/rll-url-shortener/pkg/logger"
)

type exampleList struct {
	repo repositories.Example
}

func NewExampleList(repo repositories.Example) contract.UseCase {
	return &exampleList{repo: repo}
}

// Serve partner list data
func (u *exampleList) Serve(data *appctx.Data) appctx.Response {

	p, e := u.repo.Find(data.Request.Context())

	if e != nil {
		logger.Error(logger.SetMessageFormat("[example-list] %v", e.Error()))

		return appctx.Response{
			Name: consts.ResponseInternalFailure,
		}
	}

	return appctx.Response{
		Name:    consts.ResponseSuccess,
		Message: "ok",
		Data:    p,
	}
}
