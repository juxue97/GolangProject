package main

import (
	"context"
	"fmt"

	pb "github.com/juxue97/common/api"
	"go.opentelemetry.io/otel/trace"
)

type telemetryMiddleware struct {
	next *stockService
}

func NewTelemetryMiddleware(next *stockService) *telemetryMiddleware {
	return &telemetryMiddleware{next: next}
}

func (s *telemetryMiddleware) CheckIfItemInStock(ctx context.Context, p []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CheckIfItemInStock: %v", p))

	return s.next.CheckIfItemInStock(ctx, p)
}

// func (s *telemetryMiddleware) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
// 	span := trace.SpanFromContext(ctx)
// 	span.AddEvent(fmt.Sprintf("GetItems: %v", ids))

// 	return s.next.GetItems(ctx, ids)
// }
