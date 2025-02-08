package main

import (
	"context"
	"fmt"

	pb "github.com/juxue97/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *telemetryMiddleware) GetItems(ctx context.Context) ([]*Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetItems"))
	return s.next.GetItems(ctx)
}

func (s *telemetryMiddleware) GetItem(ctx context.Context, id string) (*Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetItem: %v", id))
	return s.next.GetItem(ctx, id)
}

func (s *telemetryMiddleware) UpdateItem(ctx context.Context, id string, p *UpdateItemRequest) (*Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("UpdateItem: %v, payload: %v", id, p))
	return s.next.UpdateItem(ctx, id, p)
}

func (s *telemetryMiddleware) DeductStock(ctx context.Context, id string, quantity int) (*Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("DeductStock: %v, quantity: %v", id, quantity))

	return s.next.DeductStock(ctx, id, quantity)
}

func (s *telemetryMiddleware) GetOrderService(ctx context.Context, o *pb.Order) ([]*pb.Item, error) {
	return s.next.GetOrderService(ctx, o)
}

func (s *telemetryMiddleware) UpdateStock(ctx context.Context, id string, quantity int) (*Item, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("UpdateStock: %v, quantity: %v", id, quantity))
	return s.next.UpdateStock(ctx, id, quantity)
}

func (s *telemetryMiddleware) DeleteItem(ctx context.Context, id string) error {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("DeleteItem: %v", id))
	return s.next.DeleteItem(ctx, id)
}

func (s *telemetryMiddleware) CreateItem(ctx context.Context, p *CreateItemRequest) (primitive.ObjectID, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreateItem: %v", p))
	return s.next.CreateItem(ctx, p)
}
