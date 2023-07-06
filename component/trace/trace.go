package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Config struct {
	ServiceId string
	Env       string
	Jaeger    JaegerConfig
}

type JaegerConfig struct {
	Enable bool
	Url    string
}

func Configure(cfg *Config) (err error) {
	options := []tracesdk.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		// Always be sure to batch in production.
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(cfg.ServiceId),
			attribute.String("env", cfg.Env),
		))}
	if cfg.Jaeger.Enable {
		var exp tracesdk.SpanExporter
		exp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Jaeger.Url)))
		if err != nil {
			return
		}
		//options = append(options, tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(0.5))))
		options = append(options, tracesdk.WithBatcher(exp))
	} else {
		options = append(options, tracesdk.WithSampler(tracesdk.NeverSample()))
	}
	tp := tracesdk.NewTracerProvider(options...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}
