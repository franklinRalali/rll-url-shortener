package urlshortener

import (
	"net/url"

	"github.com/gorilla/mux"
	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	urlshortener "github.com/ralali/rll-url-shortener/internal/service/url_shortener"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	"github.com/ralali/rll-url-shortener/pkg/tracer"
)

type createShortUrl struct {
	shortUrlSvc urlshortener.URLShortener
}

func NewCreateShortURL(shortUrlSvc urlshortener.URLShortener) contract.UseCase {
	return &createShortUrl{shortUrlSvc: shortUrlSvc}
}

func (c *createShortUrl) Serve(data *appctx.Data) (response appctx.Response) {
	// TODO
	// this can possibly moved to consts
	logF := "[createShortUrl.Serve] %s"

	var (
		ctx    = tracer.SpanStartUseCase(data.Request.Context(), "Serve")
		userId = data.Request.Header.Get("user-id")
	)

	originUrl := mux.Vars(data.Request)["url"]

	var fl []logger.Field
	fl = append(
		fl,
		logger.Any("user-id", userId),
		logger.Any("origin_url", originUrl),
	)

	// decode the url as the url
	// is url encoded
	originUrl, err := url.QueryUnescape(originUrl)
	if err != nil {
		response.SetName(consts.ResponseInvalidURL)
		response.SetError(err)
		
		return response
	}

	// validate url to shorten
	_, err = url.ParseRequestURI(originUrl)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, logger.SetMessageFormat(logF, err.Error()), fl...)

		response.SetName(consts.ResponseInvalidURL)
		response.SetError(err)

		return response
	}

	req := urlshortener.ShortenURLReq{
		UserID: userId,
		OriginURL: originUrl,
	}

	res, err := c.shortUrlSvc.ShortenURL(ctx, req)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, logger.SetMessageFormat(logF, err.Error()), fl...)

		response.SetName(consts.ResponseInternalFailure)
		response.SetError(err)

		return response
	}

	response.SetName(consts.ResponseSuccess)
	response.SetData(res)

	return response
}
