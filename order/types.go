package main

import (
	"context"

	pb "github.com/juxue97/common/api"
)

type OrderService interface {
	createOrder(context.Context) error
	validateOrder(context.Context, *pb.CreateOrderRequest) error
}

type OrderStore interface {
	Create(context.Context) error
}
