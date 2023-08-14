package urlshortener

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/internal/presentations"
	urlshortener "github.com/ralali/rll-url-shortener/internal/service/url_shortener"
	"github.com/ralali/rll-url-shortener/internal/ucase/contract"
	mock_urlshortener "github.com/ralali/rll-url-shortener/mocks/service/url_shortener"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_createShortUrl_Serve(t *testing.T) {
	type fields struct {
		shortUrlSvcMock func(ctrl *gomock.Controller) urlshortener.URLShortener
	}
	type args struct {
		data func() *appctx.Data
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse appctx.Response
	}{
		{
			name: "on success shorten a url, return response 200 OK",
			fields: fields{
				shortUrlSvcMock: func(ctrl *gomock.Controller) urlshortener.URLShortener {
					mock := mock_urlshortener.NewMockURLShortener(ctrl)
					mock.EXPECT().ShortenURL(gomock.Any(), presentations.ShortenURLReq{
						UserID:    "ralali10",
						OriginURL: "http://foo.com",
					}).Return(presentations.ShortenURLRes{
						ID:        1,
						ShortCode: "abcd1234",
					}, nil)

					return mock
				},
			},
			args: args{
				data: func() *appctx.Data {
					req := httptest.NewRequest(http.MethodPost, "/{url}", nil)
					req = req.Clone(context.Background())
					req = mux.SetURLVars(req, map[string]string{"url": url.QueryEscape("http://foo.com")})
					header := make(http.Header)
					header.Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
					header.Set("User-ID", "ralali10")
					req.Header = header
					data := &appctx.Data{
						Request: req,
					}

					return data
				},
			},
			wantResponse: *appctx.NewResponse().SetData(presentations.ShortenURLRes{
				ID:        1,
				ShortCode: "abcd1234",
			}).SetName(consts.ResponseSuccess),
		},
		{
			name: "on fail with invalid url, return response 400 Bad Request",
			fields: fields{
				shortUrlSvcMock: func(ctrl *gomock.Controller) urlshortener.URLShortener {
					mock := mock_urlshortener.NewMockURLShortener(ctrl)
					return mock
				},
			},
			args: args{
				data: func() *appctx.Data {
					req := httptest.NewRequest(http.MethodPost, "/{url}", nil)
					req = req.Clone(context.Background())
					req = mux.SetURLVars(req, map[string]string{"url": url.QueryEscape("invalid_url")})
					header := make(http.Header)
					header.Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
					header.Set("User-ID", "ralali10")
					req.Header = header
					data := &appctx.Data{
						Request: req,
					}

					return data
				},
			},
			wantResponse: *appctx.NewResponse().SetName(consts.ResponseInvalidURL).SetError(&url.Error{
				Op:  "parse",
				URL: "invalid_url",
				Err: errors.New("invalid URI for request"),
			}),
		},
		{
			name: "on fail unescape url, return response 400 Bad Request",
			fields: fields{
				shortUrlSvcMock: func(ctrl *gomock.Controller) urlshortener.URLShortener {
					mock := mock_urlshortener.NewMockURLShortener(ctrl)
					return mock
				},
			},
			args: args{
				data: func() *appctx.Data {
					req := httptest.NewRequest(http.MethodPost, "/{url}", nil)
					req = req.Clone(context.Background())
					req = mux.SetURLVars(req, map[string]string{"url": "http%3%2%2Ffoo.com"})
					header := make(http.Header)
					header.Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
					header.Set("User-ID", "ralali10")
					req.Header = header
					data := &appctx.Data{
						Request: req,
					}

					return data
				},
			},
			wantResponse: *appctx.NewResponse().SetName(consts.ResponseInvalidURL).SetError(url.EscapeError("%3%")),
		},
		{
			name: "on fail from shorten url service, return 500 Internal Server Error",
			fields: fields{
				shortUrlSvcMock: func(ctrl *gomock.Controller) urlshortener.URLShortener {
					mock := mock_urlshortener.NewMockURLShortener(ctrl)
					mock.EXPECT().ShortenURL(gomock.Any(), presentations.ShortenURLReq{
						UserID:    "ralali10",
						OriginURL: "http://foo.com",
					}).Return(presentations.ShortenURLRes{}, errors.New("error"))

					return mock
				},
			},
			args: args{
				data: func() *appctx.Data {
					req := httptest.NewRequest(http.MethodPost, "/{url}", nil)
					req = req.Clone(context.Background())
					req = mux.SetURLVars(req, map[string]string{"url": url.QueryEscape("http://foo.com")})
					header := make(http.Header)
					header.Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
					header.Set("User-ID", "ralali10")
					req.Header = header
					data := &appctx.Data{
						Request: req,
					}

					return data
				},
			},
			wantResponse: *appctx.NewResponse().SetName(consts.ResponseInternalFailure).SetError(errors.New("error")),
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		shortUrlSvcMock := tt.fields.shortUrlSvcMock(mockCtrl)
		uc := createShortUrl{shortUrlSvc: shortUrlSvcMock}
		res := uc.Serve(tt.args.data())
		assert.Equal(t, tt.wantResponse, res)
	}
}

func TestNewCreateShortURL(t *testing.T) {
	type args struct {
		shortUrlSvc urlshortener.URLShortener
	}
	tests := []struct {
		name        string
		args        args
		wantUseCase contract.UseCase
	}{
		{
			name: "initiating createShortUrl use case",
			args: args{
				shortUrlSvc: urlshortener.NewURLShortener(nil, nil, appctx.Config{}, nil),
			},
			wantUseCase: &createShortUrl{
				shortUrlSvc: urlshortener.NewURLShortener(nil, nil, appctx.Config{}, nil),
			},
		},
	}
	for _, tt := range tests {
		uc := NewCreateShortURL(tt.args.shortUrlSvc)
		assert.Equal(t, tt.wantUseCase, uc)
	}
}
