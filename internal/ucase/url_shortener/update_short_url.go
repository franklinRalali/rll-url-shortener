package urlshortener

import (
	"errors"

	"github.com/gorilla/mux"
	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/internal/presentations"
	urlshortener "github.com/ralali/rll-url-shortener/internal/service/url_shortener"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/tracer"
)

type updateShortUrl struct {
	shortUrlSvc urlshortener.URLShortener
}

func NewUpdateShortURL(shortUrlSvc urlshortener.URLShortener) contract.UseCase {
	return &updateShortUrl{shortUrlSvc: shortUrlSvc}
}

func (u *updateShortUrl) Serve(data *appctx.Data) (response appctx.Response) {
	logF := "[updateShortUrl.Serve] %s"

	var (
		ctx       = tracer.SpanStartUseCase(data.Request.Context(), "Serve")
		shortCode = mux.Vars(data.Request)["short_code"]
		req       presentations.ShortURLUpdateReq
	)

	if err := data.Cast(&req); err != nil {
		response.SetName(consts.ResponseInvalidURL)
		return response
	}

	var lf []logger.Field
	lf = append(
		lf,
		logger.Any("short_code", shortCode),
		logger.Any("short_url_update_req", req),
	)

	if err := u.shortUrlSvc.UpdateShortURL(ctx, shortCode, req); err != nil {
		if errors.As(err, &urlshortener.ErrorShortURLNotFound{}) {
			response.SetName(consts.ResponseShortURLNotFound)
			return response
		}

		logger.ErrorWithContext(ctx, logger.SetMessageFormat(logF, err.Error()), lf...)
		response.SetName(consts.ResponseInternalFailure)
		return response
	}

	response.SetName(consts.ResponseSuccess)
	return response
}
