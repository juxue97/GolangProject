package main

import (
	"context"

	pb "github.com/juxue97/common/api"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, o *pb.Order) (string, error)
}
