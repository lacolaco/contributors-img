package tracing

import (
	"context"
	"log"

	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/env"
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var configuredTracer trace.Tracer

func Tracer() trace.Tracer {
	if configuredTracer == nil {
		return otel.Tracer("")
	}
	return configuredTracer
}

func installTraceProvider(cfg *config.Config) *sdktrace.TracerProvider {
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

func installPropagators() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			gcppropagator.CloudTraceOneWayPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))
}

func InitTraceProvider(cfg *config.Config) func() {
	tp := installTraceProvider(cfg)
	installPropagators()
	configuredTracer = otel.Tracer("")
	return func() {
		tp.Shutdown(context.Background())
	}
}
