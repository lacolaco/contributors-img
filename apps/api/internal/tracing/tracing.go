package tracing

import (
	"log"

	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/env"
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var DefaultTracer = otel.Tracer("")

func InitTraceProvider(cfg *config.Config) *sdktrace.TracerProvider {
	exporter, err := cloudtrace.New()
	if err != nil {
		log.Fatal(err)
	}
	tpOpts := []sdktrace.TracerProviderOption{sdktrace.WithBatcher(exporter)}
	if cfg.Env == env.EnvDevelopment {
		tpOpts = append(tpOpts, sdktrace.WithSampler(sdktrace.AlwaysSample()))
	}
	tp := sdktrace.NewTracerProvider(tpOpts...)
	otel.SetTracerProvider(tp)
	return tp
}
