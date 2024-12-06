package tracing

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// jaeger 是一个 exporter，需要传入一个 tracer 来启动
var tracer = otel.Tracer("default_tracer")

// JaegerProvider 就是 opentelemetry
// Jaeger 可以理解为收集 opentelemetry 提供的数据的后端
func InitJaegerProvider(jaegerURL, serviceName string) (func(ctx context.Context) error, error) {
	if jaegerURL == "" {
		panic("jaeger url can't be empty")
	}
	tracer = otel.Tracer(serviceName)
	// New exporter (exp)
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
	if err != nil {
		return nil, err
	}
	// New exporter 的 provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	p := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		b3Propagator,
	)
	otel.SetTextMapPropagator(p)
	return tp.Shutdown, nil
}

// name refers to span name
func Start(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

func TraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	return spanCtx.TraceID().String()
}
