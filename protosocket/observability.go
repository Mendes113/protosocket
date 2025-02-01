package protosocket

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	tracer    trace.Tracer
	meter     metric.Meter
	histogram metric.Float64Histogram
	counter   metric.Int64Counter
}

func (p *Peer) EnableTelemetry(ctx context.Context) {
	// Configuração do OpenTelemetry
	p.telemetry = &Telemetry{
		tracer: otel.GetTracerProvider().Tracer("protosocket"),
		meter:  otel.GetMeterProvider().Meter("protosocket"),
	}

	// Métricas importantes
	var err error
	p.telemetry.histogram, err = p.telemetry.meter.Float64Histogram("latency")
	if err != nil {
		return
	}
	p.telemetry.counter, err = p.telemetry.meter.Int64Counter("messages")
	if err != nil {
		return
	}
}
