package tracing

import (
	"context"
	"fmt"
	"log"

	"contrib.rocks/apps/api/internal/config"
	"contrib.rocks/libs/go/env"
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel"
	ocbridge "go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/option"
)

func Tracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("")
}

func installTraceProvider(cfg *config.Config) *sdktrace.TracerProvider {
	// Disable telemetry by Cloud Trace itself.
	// https://github.com/open-telemetry/opentelemetry-go/issues/1928#issuecomment-843644237
	exporter, err := cloudtrace.New(
		cloudtrace.WithTraceClientOptions([]option.ClientOption{option.WithTelemetryDisabled()}),
	)
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

func installOpenCensusBridge(tp *sdktrace.TracerProvider) {
	ocbridge.InstallTraceBridge(ocbridge.WithTracerProvider(tp))
}

func InitTraceProvider(cfg *config.Config) func() {
	tp := installTraceProvider(cfg)
	installPropagators()
	installOpenCensusBridge(tp)
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
