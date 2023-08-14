// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/url_shortener/contract.go

// Package mock_urlshortener is a generated GoMock package.
package mock_urlshortener

import (
	context "context"
	reflect "reflect"

	presentations "github.com/ralali/rll-url-shortener/internal/presentations"
	gomock "go.uber.org/mock/gomock"
)

// MockURLShortener is a mock of URLShortener interface.
type MockURLShortener struct {
	ctrl     *gomock.Controller
	recorder *MockURLShortenerMockRecorder
}

// MockURLShortenerMockRecorder is the mock recorder for MockURLShortener.
type MockURLShortenerMockRecorder struct {
	mock *MockURLShortener
}

// NewMockURLShortener creates a new mock instance.
func NewMockURLShortener(ctrl *gomock.Controller) *MockURLShortener {
	mock := &MockURLShortener{ctrl: ctrl}
	mock.recorder = &MockURLShortenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLShortener) EXPECT() *MockURLShortenerMockRecorder {
	return m.recorder
}

// AddVisitCount mocks base method.
func (m *MockURLShortener) AddVisitCount(ctx context.Context, shortCode string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddVisitCount", ctx, shortCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddVisitCount indicates an expected call of AddVisitCount.
func (mr *MockURLShortenerMockRecorder) AddVisitCount(ctx, shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddVisitCount", reflect.TypeOf((*MockURLShortener)(nil).AddVisitCount), ctx, shortCode)
}

// DeleteShortURLByShortCode mocks base method.
func (m *MockURLShortener) DeleteShortURLByShortCode(ctx context.Context, shortCode string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShortURLByShortCode", ctx, shortCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShortURLByShortCode indicates an expected call of DeleteShortURLByShortCode.
func (mr *MockURLShortenerMockRecorder) DeleteShortURLByShortCode(ctx, shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShortURLByShortCode", reflect.TypeOf((*MockURLShortener)(nil).DeleteShortURLByShortCode), ctx, shortCode)
}

// GetShortURL mocks base method.
func (m *MockURLShortener) GetShortURL(ctx context.Context, shortCode string) (presentations.ShortURLRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURL", ctx, shortCode)
	ret0, _ := ret[0].(presentations.ShortURLRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURL indicates an expected call of GetShortURL.
func (mr *MockURLShortenerMockRecorder) GetShortURL(ctx, shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURL", reflect.TypeOf((*MockURLShortener)(nil).GetShortURL), ctx, shortCode)
}

// GetShortURLStats mocks base method.
func (m *MockURLShortener) GetShortURLStats(ctx context.Context, shortCode string) (presentations.StatisticsRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURLStats", ctx, shortCode)
	ret0, _ := ret[0].(presentations.StatisticsRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURLStats indicates an expected call of GetShortURLStats.
func (mr *MockURLShortenerMockRecorder) GetShortURLStats(ctx, shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURLStats", reflect.TypeOf((*MockURLShortener)(nil).GetShortURLStats), ctx, shortCode)
}

// ShortenURL mocks base method.
func (m *MockURLShortener) ShortenURL(ctx context.Context, req presentations.ShortenURLReq) (presentations.ShortenURLRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShortenURL", ctx, req)
	ret0, _ := ret[0].(presentations.ShortenURLRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShortenURL indicates an expected call of ShortenURL.
func (mr *MockURLShortenerMockRecorder) ShortenURL(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShortenURL", reflect.TypeOf((*MockURLShortener)(nil).ShortenURL), ctx, req)
}

// UpdateShortURL mocks base method.
func (m *MockURLShortener) UpdateShortURL(ctx context.Context, shortCode string, req presentations.ShortURLUpdateReq) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShortURL", ctx, shortCode, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateShortURL indicates an expected call of UpdateShortURL.
func (mr *MockURLShortenerMockRecorder) UpdateShortURL(ctx, shortCode, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShortURL", reflect.TypeOf((*MockURLShortener)(nil).UpdateShortURL), ctx, shortCode, req)
}
