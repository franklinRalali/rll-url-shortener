// Package ucase
package ucase

import (
	"testing"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Serve(t *testing.T) {
	svc := NewHealthCheck()

	t.Run("test health check", func(t *testing.T) {
		result := svc.Serve(&appctx.Data{})

		assert.Equal(t, appctx.Response{
			Name:    consts.ResponseSuccess,
			Message: "ok",
		}, result)
	})
}
