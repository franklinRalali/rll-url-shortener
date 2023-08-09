// Package router
package router

import (
	"net/http"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"
	"github.com/ralali/rll-url-shortener/pkg/routerkit"
)

// httpHandlerFunc is a contract http handler for router
type httpHandlerFunc func(request *http.Request, svc contract.UseCase, conf *appctx.Config) appctx.Response

// Router is a contract router and must implement this interface
type Router interface {
	Route() *routerkit.Router
}
