// Package middleware
package middleware

import (
	"net/http"
	"strings"
	"fmt"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/consts"
	"github.com/ralali/rll-url-shortener/pkg/logger"
)

// ValidateContentType header
func ValidateContentType(r *http.Request, conf *appctx.Config) string {

	if ct := strings.ToLower(r.Header.Get(`Content-Type`)) ; ct != `application/json` {
		logger.Warn(fmt.Sprintf("[middleware] invalid content-type %s", ct ))

		return consts.ResponseUnprocessableEntity
	}


	return consts.MiddlewarePassed
}
