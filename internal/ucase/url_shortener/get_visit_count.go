package urlshortener

import (
	"errors"

	"github.com/gorilla/mux"
	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	urlshortener "github.com/ralali/rll-url-shortener/internal/service/url_shortener"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/tracer"
)

type getVisitCount struct {
	shortUrlSvc urlshortener.URLShortener
}

func NewGetVisitCount(shortUrlSvc urlshortener.URLShortener) contract.UseCase {
	return &getVisitCount{shortUrlSvc: shortUrlSvc}
}

func (g *getVisitCount) Serve(data *appctx.Data) (response appctx.Response) {
	logF := "[getVisitCount.Serve] %s"

	var (
		ctx       = tracer.SpanStartUseCase(data.Request.Context(), "Serve")
		shortCode = mux.Vars(data.Request)["short_code"]
	)

	var lf []logger.Field
	lf = append(lf, logger.Any("short_code", shortCode))

	res, err := g.shortUrlSvc.GetShortURLStats(ctx, shortCode)
	if err != nil {
		if errors.As(err, new(urlshortener.ErrorShortURLNotFound)) {
			response.SetName(consts.ResponseShortURLNotFound)
			return response
		}

		logger.ErrorWithContext(ctx, logger.SetMessageFormat(logF, err.Error()), lf...)

		response.SetName(consts.ResponseInternalFailure)
		return response
	}

	response.SetName(consts.ResponseSuccess)
	response.SetData(res)
	return response
}
