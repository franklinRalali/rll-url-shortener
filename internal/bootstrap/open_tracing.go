package bootstrap

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/ralali/rll-url-shortener/internal/appctx"
	"github.com/ralali/rll-url-shortener/pkg/logger"
	tracerCfg "github.com/ralali/rll-url-shortener/pkg/tracer"
)

// RegistryOpenTracing setup
func RegistryOpenTracing(cfg *appctx.Config) opentracing.Tracer {

	if !cfg.APM.Enable {
		return opentracing.NoopTracer{}
	}

	lf := logger.NewFields(logger.EventName("TracerInitiated"))
	logger.Debug(fmt.Sprint("apm address : ", cfg.APM.Address), lf...)
	tr := opentracer.New(
		tracer.WithAgentAddr(cfg.APM.Address),
		tracer.WithService(cfg.APM.Name),
		tracer.WithGlobalTag("env", cfg.App.Env),
	)

	tracerCfg.New(cfg.App.AppName)
	opentracing.SetGlobalTracer(tr)
	return tr

}
