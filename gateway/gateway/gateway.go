package gateway

import (
	"context"

	pb "github.com/juxue97/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(context.Context, string, string) (*pb.Order, error)
}
