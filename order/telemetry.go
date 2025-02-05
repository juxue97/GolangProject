package main

import (
	"context"
	"fmt"

	pb "github.com/juxue97/common/api"
	"go.opentelemetry.io/otel/trace"
)

type telemetryMiddleware struct {
	next OrderService
}

func NewtelemetryMiddleware(next OrderService) *telemetryMiddleware {
	return &telemetryMiddleware{next: next}
}

func (s *telemetryMiddleware) createOrder(ctx context.Context, payload *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreateOrder: %v, items: %v", payload, items))

	return s.next.createOrder(ctx, payload, items)
}

func (s *telemetryMiddleware) getOrder(ctx context.Context, payload *pb.GetOrderRequest) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetOrder: %v", payload))

	return s.next.getOrder(ctx, payload)
}

func (s *telemetryMiddleware) validateOrder(ctx context.Context, payload *pb.CreateOrderRequest) ([]*pb.Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("ValidateOrder: %v", payload))
	return s.next.validateOrder(ctx, payload)
}

func (s *telemetryMiddleware) updateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("UpdateOrder: %v", o))
	return s.next.updateOrder(ctx, o)
}
