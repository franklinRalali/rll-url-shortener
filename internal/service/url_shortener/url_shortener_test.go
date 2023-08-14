package urlshortener

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/entity"
	"github.com/ralali/rll-url-shortener/internal/presentations"
	"github.com/ralali/rll-url-shortener/internal/repositories"
	mock_cache "github.com/ralali/rll-url-shortener/mocks/cache"
	mock_rands "github.com/ralali/rll-url-shortener/mocks/rands"
	mock_repositories "github.com/ralali/rll-url-shortener/mocks/repositories"
	"github.com/ralali/rll-url-shortener/pkg/cache"
	"github.com/ralali/rll-url-shortener/pkg/rands"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_urlShortener_ShortenURL(t *testing.T) {
	type fields struct {
		urlsRepoMock func(ctrl *gomock.Controller) repositories.URLs
		cacherMock   func(ctrl *gomock.Controller) cache.Cacher
		conf         appctx.Config
		randGenMock  func(ctrl *gomock.Controller) rands.Randoms
	}
	type args struct {
		ctx context.Context
		req presentations.ShortenURLReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes presentations.ShortenURLRes
		wantErr bool
	}{
		{
			name: "on success shortening a url, returns no error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().Upsert(context.Background(), entity.URL{
						URL:       "http://foo.com",
						ShortCode: "abcd1234",
					}).Return(uint64(1), nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Set(context.Background(), "short-url-origin-abcd1234", "http://foo.com", 5*time.Minute).Return(nil)

					return mock
				},
				conf: appctx.Config{
					App: &appctx.Common{
						ShortURLHost: "http://short.com",
					},
				},
				randGenMock: func(ctrl *gomock.Controller) rands.Randoms {
					mock := mock_rands.NewMockRandoms(ctrl)
					mock.EXPECT().String(8, rands.DefaultCharSet).Return("abcd1234", nil)

					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: presentations.ShortenURLReq{
					OriginURL: "http://foo.com",
				},
			},
			wantRes: presentations.ShortenURLRes{
				ID:        1,
				ShortCode: "abcd1234",
			},
		},
		{
			name: "on fail random generator fails, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
				randGenMock: func(ctrl *gomock.Controller) rands.Randoms {
					mock := mock_rands.NewMockRandoms(ctrl)
					mock.EXPECT().String(8, rands.DefaultCharSet).Return("", errors.New("error"))

					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: presentations.ShortenURLReq{
					OriginURL: "https://foo.com",
				},
			},
			wantRes: presentations.ShortenURLRes{},
			wantErr: true,
		},
		{
			name: "on fail insert to database, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().Upsert(context.Background(), entity.URL{
						URL:       "http://foo.com",
						ShortCode: "abcd1234",
					}).Return(uint64(0), errors.New("error"))

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
				randGenMock: func(ctrl *gomock.Controller) rands.Randoms {
					mock := mock_rands.NewMockRandoms(ctrl)
					mock.EXPECT().String(8, rands.DefaultCharSet).Return("abcd1234", nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: presentations.ShortenURLReq{
					OriginURL: "http://foo.com",
				},
			},
			wantRes: presentations.ShortenURLRes{},
			wantErr: true,
		},
		{
			name: "on fail insert to cache, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().Upsert(context.Background(), entity.URL{
						URL:       "http://foo.com",
						ShortCode: "abcd1234",
					}).Return(uint64(1), nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Set(context.Background(), "short-url-origin-abcd1234", "http://foo.com", originUrlCacheTimeout).Return(errors.New("error"))
					return mock
				},
				randGenMock: func(ctrl *gomock.Controller) rands.Randoms {
					mock := mock_rands.NewMockRandoms(ctrl)
					mock.EXPECT().String(8, rands.DefaultCharSet).Return("abcd1234", nil)
					return mock
				},
			},
			args: args{
				ctx: context.Background(),
				req: presentations.ShortenURLReq{
					OriginURL: "http://foo.com",
				},
			},
			wantRes: presentations.ShortenURLRes{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		repoMock := tt.fields.urlsRepoMock(mockCtrl)
		cacheMock := tt.fields.cacherMock(mockCtrl)
		randGenMock := tt.fields.randGenMock(mockCtrl)
		svc := NewURLShortener(repoMock, cacheMock, tt.fields.conf, randGenMock)
		got, err := svc.ShortenURL(tt.args.ctx, tt.args.req)

		assert.Equal(t, tt.wantRes, got)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func Test_urlShortener_GetShortURL(t *testing.T) {
	type fields struct {
		urlsRepoMock func(ctrl *gomock.Controller) repositories.URLs
		cacherMock   func(ctrl *gomock.Controller) cache.Cacher
		conf         appctx.Config
	}
	type args struct {
		ctx       context.Context
		shortCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes presentations.ShortURLRes
		wantErr bool
	}{
		{
			name: "on success get short url from cache, returns no error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Get(context.Background(), "short-url-origin-abcd1234").Return([]byte("http://foo.com"), nil)

					return mock
				},
				conf: appctx.Config{
					App: &appctx.Common{
						ShortURLHost: "http://short.com",
					},
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantRes: presentations.ShortURLRes{
				OriginURL: "http://foo.com",
				ShortURL:  "http://short.com/abcd1234",
			},
			wantErr: false,
		},
		{
			name: "on success get short url from database, returns no error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("http://foo.com", nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Get(context.Background(), "short-url-origin-abcd1234").Return(nil, nil)
					mock.EXPECT().Set(context.Background(), "short-url-origin-abcd1234", "http://foo.com", originUrlCacheTimeout).Return(nil)

					return mock
				},
				conf: appctx.Config{
					App: &appctx.Common{
						ShortURLHost: "http://short.com",
					},
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantRes: presentations.ShortURLRes{
				OriginURL: "http://foo.com",
				ShortURL:  "http://short.com/abcd1234",
			},
			wantErr: false,
		},
		{
			name: "on fail fetch origin url from cache, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					return mock
				},

				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Get(context.Background(), "short-url-origin-abcd1234").Return(nil, errors.New("error"))

					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantRes: presentations.ShortURLRes{},
			wantErr: true,
		},
		{
			name: "on fail fetch origin url from database, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("", errors.New("error"))

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Get(context.Background(), "short-url-origin-abcd1234").Return(nil, nil)

					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantRes: presentations.ShortURLRes{},
			wantErr: true,
		},
		{
			name: "on fail origin url is not exist in database, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("", ErrorShortURLNotFound{ShortCode: "abcd1234"})

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Get(context.Background(), "short-url-origin-abcd1234").Return(nil, nil)

					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantRes: presentations.ShortURLRes{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		urlsRepoMock := tt.fields.urlsRepoMock(mockCtrl)
		cacheMock := tt.fields.cacherMock(mockCtrl)
		svc := NewURLShortener(urlsRepoMock, cacheMock, tt.fields.conf, nil)
		got, err := svc.GetShortURL(tt.args.ctx, tt.args.shortCode)
		assert.Equal(t, tt.wantRes, got)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func Test_urlShortener_UpdateShortURL(t *testing.T) {
	type fields struct {
		urlsRepoMock func(ctrl *gomock.Controller) repositories.URLs
		cacherMock   func(ctrl *gomock.Controller) cache.Cacher
		conf         appctx.Config
	}
	type args struct {
		ctx       context.Context
		shortCode string
		req       presentations.ShortURLUpdateReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "on success updating short url, returns no error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("http://oldfoo.com", nil)
					mock.EXPECT().UpdateByShortCode(context.Background(), "abcd1234", entity.URL{
						URL: "http://foo.com",
					}).Return(nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Set(context.Background(), "short-url-origin-abcd1234", "http://foo.com", originUrlCacheTimeout).Return(nil)

					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
				req: presentations.ShortURLUpdateReq{
					OriginURL: "http://foo.com",
				},
			},
			wantErr: false,
		},
		{
			name: "on fail checking short url exist, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("", errors.New("error"))
					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
				req: presentations.ShortURLUpdateReq{
					OriginURL: "http://foo.com",
				},
			},
			wantErr: true,
		},
		{
			name: "on fail short url is not exist, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("", sql.ErrNoRows)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
				req: presentations.ShortURLUpdateReq{
					OriginURL: "http://foo.com",
				},
			},
			wantErr: true,
		},
		{
			name: "on fail updating origin url to database, return error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("http://oldfoo.com", nil)
					mock.EXPECT().UpdateByShortCode(context.Background(), "abcd1234", entity.URL{
						URL: "http://foo.com",
					}).Return(errors.New("error"))

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
				req: presentations.ShortURLUpdateReq{
					OriginURL: "http://foo.com",
				},
			},
			wantErr: true,
		},
		{
			name: "on fail updating origin url to cache, return error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("http://oldfoo.com", nil)
					mock.EXPECT().UpdateByShortCode(context.Background(), "abcd1234", entity.URL{
						URL: "http://foo.com",
					}).Return(nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Set(context.Background(), "short-url-origin-abcd1234", "http://foo.com", originUrlCacheTimeout).Return(errors.New("error"))

					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
				req: presentations.ShortURLUpdateReq{
					OriginURL: "http://foo.com",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		urlsRepoMock := tt.fields.urlsRepoMock(mockCtrl)
		cacheMock := tt.fields.cacherMock(mockCtrl)
		svc := NewURLShortener(urlsRepoMock, cacheMock, tt.fields.conf, nil)
		err := svc.UpdateShortURL(tt.args.ctx, tt.args.shortCode, tt.args.req)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func Test_urlShortener_DeleteShortURLByShortCode(t *testing.T) {
	type fields struct {
		urlsRepoMock func(ctrl *gomock.Controller) repositories.URLs
		cacherMock   func(ctrl *gomock.Controller) cache.Cacher
		conf         appctx.Config
	}
	type args struct {
		ctx       context.Context
		shortCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "on success deleting short url, returns no error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("http://foo.com", nil)
					mock.EXPECT().DeleteByShortCode(context.Background(), "abcd1234").Return(nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Delete(context.Background(), "short-url-origin-abcd1234").Return(nil)

					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantErr: false,
		},
		{
			name: "on fail checking origin url exist, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("", errors.New("error"))

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantErr: true,
		},
		{
			name: "on fail short url is not exist, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("", sql.ErrNoRows)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantErr: true,
		},
		{
			name: "on fail deleting short url from database, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("http://foo.com", nil)
					mock.EXPECT().DeleteByShortCode(context.Background(), "abcd1234").Return(errors.New("error"))

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantErr: true,
		},
		{
			name: "on fail deleting origin url from cache, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneOriginURLByShortCode(context.Background(), "abcd1234").Return("http://foo.com", nil)
					mock.EXPECT().DeleteByShortCode(context.Background(), "abcd1234").Return(nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					mock.EXPECT().Delete(context.Background(), "short-url-origin-abcd1234").Return(errors.New("error"))

					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		urlsRepoMock := tt.fields.urlsRepoMock(mockCtrl)
		cacheMock := tt.fields.cacherMock(mockCtrl)
		svc := NewURLShortener(urlsRepoMock, cacheMock, tt.fields.conf, nil)
		err := svc.DeleteShortURLByShortCode(tt.args.ctx, tt.args.shortCode)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func Test_urlShortener_GetShortURLStats(t *testing.T) {
	type fields struct {
		urlsRepoMock func(ctrl *gomock.Controller) repositories.URLs
		cacherMock   func(ctrl *gomock.Controller) cache.Cacher
		conf         appctx.Config
	}
	type args struct {
		ctx       context.Context
		shortCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes presentations.StatisticsRes
		wantErr bool
	}{
		{
			name: "on success get short url statistics, returns no error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneVisitCountByShortCode(context.Background(), "abcd1234").Return(uint64(1), nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantRes: presentations.StatisticsRes{
				VisitCount: 1,
			},
			wantErr: false,
		},
		{
			name: "on fail fetching visit count, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneVisitCountByShortCode(context.Background(), "abcd1234").Return(uint64(0), errors.New("error"))

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantRes: presentations.StatisticsRes{},
			wantErr: true,
		},
		{
			name: "on fail short url is not exist in database, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().FindOneVisitCountByShortCode(context.Background(), "abcd1234").Return(uint64(0), sql.ErrNoRows)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantRes: presentations.StatisticsRes{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		urlsRepoMock := tt.fields.urlsRepoMock(mockCtrl)
		cacheMock := tt.fields.cacherMock(mockCtrl)
		svc := NewURLShortener(urlsRepoMock, cacheMock, tt.fields.conf, nil)
		got, err := svc.GetShortURLStats(tt.args.ctx, tt.args.shortCode)
		assert.Equal(t, tt.wantRes, got)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func Test_urlShortener_AddVisitCount(t *testing.T) {
	type fields struct {
		urlsRepoMock func(ctrl *gomock.Controller) repositories.URLs
		cacherMock   func(ctrl *gomock.Controller) cache.Cacher
		conf         appctx.Config
	}
	type args struct {
		ctx       context.Context
		shortCode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "on success add visit count, returns no error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().AddVisitCountByShortCode(context.Background(), "abcd1234", uint(1)).Return(nil)

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantErr: false,
		},
		{
			name: "on fail add visit account, returns error",
			fields: fields{
				urlsRepoMock: func(ctrl *gomock.Controller) repositories.URLs {
					mock := mock_repositories.NewMockURLs(ctrl)
					mock.EXPECT().AddVisitCountByShortCode(context.Background(), "abcd1234", uint(1)).Return(errors.New("error"))

					return mock
				},
				cacherMock: func(ctrl *gomock.Controller) cache.Cacher {
					mock := mock_cache.NewMockCacher(ctrl)
					return mock
				},
			},
			args: args{
				ctx:       context.Background(),
				shortCode: "abcd1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		urlsRepoMock := tt.fields.urlsRepoMock(mockCtrl)
		cacheMock := tt.fields.cacherMock(mockCtrl)
		svc := NewURLShortener(urlsRepoMock, cacheMock, tt.fields.conf, nil)
		err := svc.AddVisitCount(tt.args.ctx, tt.args.shortCode)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestNewURLShortener(t *testing.T) {
	type args struct {
		urlsRepo repositories.URLs
		cacher   cache.Cacher
		conf     appctx.Config
		randGen  rands.Randoms
	}
	tests := []struct {
		name    string
		args    args
		wantSvc URLShortener
	}{
		{
			name: "initiating implementation of URLShortener",
			args: args{
				urlsRepo: repositories.NewUrls(nil),
				cacher:   cache.NewCache(nil),
				conf:     appctx.Config{},
				randGen:  rands.New(nil),
			},
			wantSvc: &urlShortener{
				urlsRepo: repositories.NewUrls(nil),
				cacher:   cache.NewCache(nil),
				conf:     appctx.Config{},
				randGen:  rands.New(nil),
			},
		},
	}
	for _, tt := range tests {
		svc := NewURLShortener(tt.args.urlsRepo, tt.args.cacher, tt.args.conf, tt.args.randGen)
		assert.Equal(t, tt.wantSvc, svc)
	}
}

func TestErrorShortURLNotFound_Error(t *testing.T) {
	type fields struct {
		ShortCode string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "initiating ErrShortURLNotFound",
			fields: fields{
				ShortCode: "abcd1234",
			},
			want: "short URL with short code abcd1234 is not found",
		},
	}
	for _, tt := range tests {
		err := ErrorShortURLNotFound{ShortCode: tt.fields.ShortCode}
		assert.Equal(t, tt.want, err.Error())
	}
}
