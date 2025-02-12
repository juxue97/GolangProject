package gateway

import (
	"context"
	"log"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/common/discovery"
)

var orderServiceName = "orders"

type ordersGateway struct {
	registry discovery.Registry
}

func NewOrdersGateway(registry discovery.Registry) *ordersGateway {
	return &ordersGateway{registry: registry}
}

func (g *ordersGateway) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
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

func (g *ordersGateway) GetOrder(ctx context.Context, orderID string, customerID string) (*pb.Order, error) {
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
