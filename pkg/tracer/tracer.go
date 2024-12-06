// Package tracer реализует провайдер и экспортер трассировки при помощи OpenTelemetry
package tracer

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Tracer работа с трассировкой при помощи OpenTelemetry
type Tracer struct {
	exporter sdktrace.SpanExporter
	provider *sdktrace.TracerProvider
}

// NewTracer инициализация трассировки
func NewTracer(ctx context.Context, instance string, endpoint string) (*Tracer, error) {
	exporter, provider, err := initOtel(ctx, instance, endpoint)
	if err != nil {
		return nil, err
	}
	return &Tracer{
		exporter: exporter,
		provider: provider,
	}, nil
}

// Shutdown shuts down the trace exporter and trace provider.
func (t *Tracer) Shutdown(ctx context.Context) error {

	// Shutdown the trace provider.
	err := t.provider.Shutdown(ctx)

	// Shutdown the trace exporter.
	if err1 := t.exporter.Shutdown(ctx); err1 != nil {
		err = errors.Join(err, err1)
	}

	if err != nil {
		return err
	}
	return nil
}

// initOtel init configures an OpenTelemetry exporter and trace provider.
func initOtel(ctx context.Context, instance string, endpoint string) (sdktrace.SpanExporter, *sdktrace.TracerProvider, error) {

	exporter, err := otlptracegrpc.New( // grpc экспортер
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(instance),
			),
		),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(
			sdktrace.AlwaysSample(),
			//trace.ParentBased(trace.TraceIDRatioBased(0.2)), // если нет родительского семплера, то 20% сэмплируем, иначе используем родительский
			//trace.NeverSample(),
		),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return exporter, tracerProvider, nil
}
