package trace

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

func TraceID(ctx context.Context) (res string) {
	if ctx == nil {
		return
	}
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		res = span.TraceID().String()
	}
	return
}
func SpanID(ctx context.Context) (res string) {
	if ctx == nil {
		return
	}
	if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
		res = span.SpanID().String()
	}
	return
}
