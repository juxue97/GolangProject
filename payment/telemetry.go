package main

import (
	"context"
	"fmt"

	pb "github.com/juxue97/common/api"
	"go.opentelemetry.io/otel/trace"
)

type telemetryMiddleware struct {
	next PaymentService
}

func NewtelemetryMiddleware(next PaymentService) *telemetryMiddleware {
	return &telemetryMiddleware{next: next}
}

func (s *telemetryMiddleware) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf(
		"CreatePayment: %v", o,
	))

	return s.next.CreatePayment(ctx, o)
}
