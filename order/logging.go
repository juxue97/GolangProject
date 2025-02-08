package main

import (
	"context"
	"time"

	pb "github.com/juxue97/common/api"
	"go.uber.org/zap"
)

type loggingMiddleware struct {
	next *telemetryMiddleware
}

func NewloggingMiddleware(next *telemetryMiddleware) *loggingMiddleware {
	return &loggingMiddleware{next: next}
}

func (s *loggingMiddleware) createOrder(ctx context.Context, payload *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CreateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.createOrder(ctx, payload, items)
}

func (s *loggingMiddleware) getOrder(ctx context.Context, payload *pb.GetOrderRequest) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("GetOrder", zap.Duration("took", time.Since(start)))
	}()
	return s.next.getOrder(ctx, payload)
}

func (s *loggingMiddleware) getOrderForStock(ctx context.Context, payload *pb.GetOrderRequest) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("GetOrderForStock", zap.Duration("took", time.Since(start)))
	}()
	return s.next.getOrderForStock(ctx, payload)
}

func (s *loggingMiddleware) validateOrder(ctx context.Context, payload *pb.CreateOrderRequest) ([]*pb.Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("ValidateOrder", zap.Duration("took", time.Since(start)))
	}()
	return s.next.validateOrder(ctx, payload)
}

func (s *loggingMiddleware) updateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("UpdateOrder", zap.Duration("took", time.Since(start)))
	}()
	return s.next.updateOrder(ctx, o)
}
