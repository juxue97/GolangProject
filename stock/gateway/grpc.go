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

func (g *gateway) GetOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(context.Background(), orderServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	ordersClient := pb.NewOrderServiceClient(conn)

	res, err := ordersClient.GetOrderForStockUpdate(ctx, &pb.GetOrderRequest{
		OrderID:    o.ID,
		CustomerID: o.CustomerID,
	})
	return res, err
}
