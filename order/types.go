package main

import (
	"context"

	pb "github.com/juxue97/common/api"
)

type OrderService interface {
	createOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	validateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	getOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	updateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrderStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(context.Context, string, string) (*pb.Order, error)
	Update(context.Context, string, *pb.Order) error
}
