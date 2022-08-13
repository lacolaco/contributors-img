package tracing

import (
	"context"
	"fmt"
	"log"

	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/goutils/env"
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	octrace "go.opencensus.io/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func Tracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("")
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
			gcppropagator.CloudTraceFormatPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))
}

func installOpenCensusBridge() {
	tracer := otel.GetTracerProvider().Tracer("ocbridge")
	octrace.DefaultTracer = opencensus.NewTracer(tracer)
}

func InitTraceProvider(cfg *config.Config) func() {
	tp := installTraceProvider(cfg)
	installPropagators()
	installOpenCensusBridge()
	return func() {
		tp.Shutdown(context.Background())
	}
}

func buildTraceName(projectID string, traceID string) string {
	if projectID == "" || traceID == "" {
		return ""
	}
	return fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
}
