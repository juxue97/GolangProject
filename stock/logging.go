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

func (s *loggingMiddleware) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("GetItems", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetItems(ctx, ids)
}
