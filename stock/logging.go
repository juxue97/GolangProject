package main

import (
	"context"
	"time"

	pb "github.com/juxue97/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type loggingMiddleware struct {
	next *telemetryMiddleware
}

func NewLoggingMiddleware(next *telemetryMiddleware) *loggingMiddleware {
	return &loggingMiddleware{next: next}
}

func (s *loggingMiddleware) CheckIfItemInStock(ctx context.Context, p []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CheckIfItemInStock", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CheckIfItemInStock(ctx, p)
}

func (s *loggingMiddleware) GetItems(ctx context.Context) ([]*Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("GetItems", zap.Duration("took", time.Since(start)))
	}()
	return s.next.GetItems(ctx)
}

func (s *loggingMiddleware) GetItem(ctx context.Context, id string) (*Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("GetItem", zap.Duration("took", time.Since(start)))
	}()
	return s.next.GetItem(ctx, id)
}

func (s *loggingMiddleware) UpdateItem(ctx context.Context, id string, p *UpdateItemRequest) (*Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("UpdateItem", zap.Duration("took", time.Since(start)))
	}()
	return s.next.UpdateItem(ctx, id, p)
}

func (s *loggingMiddleware) UpdateStock(ctx context.Context, id string, quantity int) (*Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("UpdateStock", zap.Duration("took", time.Since(start)))
	}()
	return s.next.UpdateStock(ctx, id, quantity)
}

func (s *loggingMiddleware) DeductStock(ctx context.Context, id string, quantity int) (*Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("DeductStock", zap.Duration("took", time.Since(start)))
	}()
	return s.next.DeductStock(ctx, id, quantity)
}

func (s *loggingMiddleware) GetOrderService(ctx context.Context, o *pb.Order) ([]*pb.Item, error) {
	return s.next.GetOrderService(ctx, o)
}

func (s *loggingMiddleware) DeleteItem(ctx context.Context, id string) error {
	start := time.Now()
	defer func() {
		zap.L().Info("DeleteItem", zap.Duration("took", time.Since(start)))
	}()
	return s.next.DeleteItem(ctx, id)
}

func (s *loggingMiddleware) CreateItem(ctx context.Context, p *CreateItemRequest) (primitive.ObjectID, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CreateItem", zap.Duration("took", time.Since(start)))
	}()
	return s.next.CreateItem(ctx, p)
}
