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

func (s *loggingMiddleware) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CreatePayment", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreatePayment(ctx, o)
}
