package urlshortener

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	urlshortener "github.com/ralali/rll-url-shortener/internal/service/url_shortener"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/tracer"
)

type getShortUrl struct {
	shortUrlSvc urlshortener.URLShortener
}

func NewGetShortURL(shortUrlSvc urlshortener.URLShortener) contract.UseCase {
	return &getShortUrl{shortUrlSvc: shortUrlSvc}
}

func (g *getShortUrl) Serve(data *appctx.Data) (response appctx.Response) {
	// TODO
	// this can possibly moved to consts
	logF := "[getShortUrl.Serve] %s"

	var (
		ctx       = tracer.SpanStartUseCase(data.Request.Context(), "Serve")
		shortCode = mux.Vars(data.Request)["short_code"]
	)

	lf := []logger.Field{
		logger.Any("short_code", shortCode),
	}

	res, err := g.shortUrlSvc.GetShortURL(ctx, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			response.SetName(consts.ResponseShortURLNotFound)
			return response
		}

		logger.ErrorWithContext(ctx, logger.SetMessageFormat(logF, err.Error()), lf...)
		response.SetName(consts.ResponseInternalFailure)
		response.SetError(err)
		return response
	}

	response.SetName(consts.ResponseSuccess)
	response.SetData(res)

	return response
}
