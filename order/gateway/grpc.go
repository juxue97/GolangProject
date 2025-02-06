package gateway

import (
	"context"
	"log"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/common/discovery"
)

var stockServiceName = "stocks"

type gateway struct {
	registry discovery.Registry
}

func NewGateway(registry discovery.Registry) *gateway {
	return &gateway{registry: registry}
}

func (g *gateway) CheckIfItemsInStock(ctx context.Context, customerID string, items []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	conn, err := discovery.ServiceConnection(context.Background(), stockServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	ordersClient := pb.NewStockServiceClient(conn)

	res, err := ordersClient.CheckIfItemsInStock(ctx, &pb.CheckIfItemsInStockRequest{
		Items: items,
	})

	return res.InStock, res.Items, err
}
