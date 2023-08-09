// Package bootstrap
package bootstrap

import (
	"github.com/thedevsaddam/govalidator"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/internal/validator"
)

// RegistryValidatorRules initialize validation rules
func RegistryValidatorRules(cfg *appctx.Config) {
	for k, f := range validator.Rules(cfg.App.Env) {
		govalidator.AddCustomRule(k, f)
	}
}
