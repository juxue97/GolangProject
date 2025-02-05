package gateway

import (
	"context"
	"log"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/common/discovery"
)

var orderServiceName = "orders"

type gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *gateway {
	return &gateway{registry: registry}
}

func (g *gateway) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(context.Background(), orderServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewOrderServiceClient(conn)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: payload.CustomerID,
		Items:      payload.Items,
	})
}

func (g *gateway) GetOrder(ctx context.Context, orderID string, customerID string) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(context.Background(), orderServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewOrderServiceClient(conn)

	return c.GetOrder(ctx, &pb.GetOrderRequest{
		OrderID:    orderID,
		CustomerID: customerID,
	})
}
