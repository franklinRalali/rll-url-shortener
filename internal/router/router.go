// Package router
package router

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/bootstrap"
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/internal/handler"
	"github.com/ralali/rll-url-shortener/internal/middleware"
	"github.com/ralali/rll-url-shortener/internal/repositories"
	urlshortenersvc "github.com/ralali/rll-url-shortener/internal/service/url_shortener"
	"github.com/ralali/rll-url-shortener/internal/ucase"
	"github.com/ralali/rll-url-shortener/pkg/cache"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/routerkit"

	//"github.com/ralali/rll-url-shortener/pkg/mariadb"
	//"github.com/ralali/rll-url-shortener/internal/repositories"
	//"github.com/ralali/rll-url-shortener/internal/ucase/example"

	ucaseContract "github.com/ralali/rll-url-shortener/internal/ucase/contract"
	urlshorteneruc "github.com/ralali/rll-url-shortener/internal/ucase/url_shortener"
)

type router struct {
	config *appctx.Config
	router *routerkit.Router
}

// NewRouter initialize new router wil return Router Interface
func NewRouter(cfg *appctx.Config) Router {
	bootstrap.RegistryValidatorRules(cfg)
	bootstrap.RegistryMessage()
	bootstrap.RegistryLogger(cfg)
	bootstrap.RegistrySnowflake()

	return &router{
		config: cfg,
		router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.AppName)),
	}
}

func (rtr *router) handle(hfn httpHandlerFunc, svc ucaseContract.UseCase, mdws ...middleware.MiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := appctx.Response{
					Name: consts.ResponseInternalFailure,
				}

				res.BuildMessage()
				logger.Error(logger.SetMessageFormat("error %v", string(debug.Stack())))
				json.NewEncoder(w).Encode(res)
				return
			}
		}()

		var st time.Time
		var lt time.Duration

		st = time.Now()

		ctx := context.WithValue(r.Context(), "access", map[string]interface{}{
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
			"method":    r.Method,
		})

		req := r.WithContext(ctx)
		defaultLang := rtr.defaultLang(req.Header.Get(consts.HeaderLanguageKey))

		if status := middleware.FilterFunc(rtr.config, req, mdws); status != consts.MiddlewarePassed {
			rtr.response(w, appctx.Response{
				Name: status,
				Lang: defaultLang,
			})

			return
		}

		resp := hfn(req, svc, rtr.config)

		resp.Lang = defaultLang

		rtr.response(w, resp)

		lt = time.Since(st)
		logger.AccessLog("access log",
			logger.Any("tag", "go-access"),
			logger.Any("http.path", req.URL.Path),
			logger.Any("http.method", req.Method),
			logger.Any("http.agent", req.UserAgent()),
			logger.Any("http.referer", req.Referer()),
			logger.Any("http.status", resp.GetCode()),
			logger.Any("http.latency", lt.Seconds()),
		)
	}
}

// response prints as a json and formatted string for DGP legacy
func (rtr *router) response(w http.ResponseWriter, resp appctx.Response) {

	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)

	defer func() {
		resp.BuildMessage()
		w.WriteHeader(resp.GetCode())
		json.NewEncoder(w).Encode(resp)
	}()

	return

}

// Route preparing http router and will return mux router object
func (rtr *router) Route() *routerkit.Router {

	root := rtr.router.PathPrefix("/").Subrouter()
	in := root.PathPrefix("/in/").Subrouter()

	// open tracer setup
	bootstrap.RegistryOpenTracing(rtr.config)

	db := bootstrap.RegistryMariaMasterSlave(rtr.config.WriteDB, rtr.config.ReadDB, time.Local.String())
	cacher := cache.NewCache(bootstrap.RegistryRedisNative(rtr.config))

	// repositories
	urlsRepo := repositories.NewUrls(db)

	// services
	urlShortSvc := urlshortenersvc.NewURLShortener(urlsRepo, cacher, *rtr.config)

	// use case
	healthy := ucase.NewHealthCheck()
	createShortUrl := urlshorteneruc.NewCreateShortURL(urlShortSvc)
	getShortUrl := urlshorteneruc.NewGetShortURL(urlShortSvc)
	updateShortUrl := urlshorteneruc.NewUpdateShortURL(urlShortSvc)
	deleteShortUrl := urlshorteneruc.NewDeleteShortURL(urlShortSvc)
	getShortUrlStats := urlshorteneruc.NewGetVisitCount(urlShortSvc)

	// healthy
	in.HandleFunc("/health", rtr.handle(
		handler.HttpRequest,
		healthy,
	)).Methods(http.MethodGet)

	// create shorten url
	in.HandleFunc("/{url}", rtr.handle(
		handler.HttpRequest,
		createShortUrl,
	)).Methods(http.MethodPost)

	// get short url
	in.HandleFunc("/{short_code}", rtr.handle(
		handler.HttpRequest,
		getShortUrl,
	)).Methods(http.MethodGet)

	// update short url
	in.HandleFunc("/{short_code}", rtr.handle(
		handler.HttpRequest,
		updateShortUrl,
	)).Methods(http.MethodPut)

	// delete short url
	in.HandleFunc("/{short_code}", rtr.handle(
		handler.HttpRequest,
		deleteShortUrl,
	)).Methods(http.MethodDelete)

	// get short url statistics
	in.HandleFunc("/{short_code}/stats", rtr.handle(
		handler.HttpRequest,
		getShortUrlStats,
	)).Methods(http.MethodGet)

	return rtr.router

}

func (rtr *router) defaultLang(l string) string {
	if len(l) == 0 {
		return rtr.config.App.DefaultLang
	}

	return l
}
