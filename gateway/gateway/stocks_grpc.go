package gateway

import (
	"context"
	"log"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/common/discovery"
)

var stockServiceName = "stocks"

type stocksGateway struct {
	registry discovery.Registry
}

func NewStocksGateway(registry discovery.Registry) *stocksGateway {
	return &stocksGateway{registry: registry}
}

func (g *stocksGateway) CreateItem(ctx context.Context, p *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), stockServiceName, nil)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewStockServiceClient(conn)

	return c.CreateStockItem(ctx, p)
}

func (g *stocksGateway) GetItems(ctx context.Context) (*pb.GetStockItemsResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), stockServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewStockServiceClient(conn)

	return c.GetStockItems(ctx, nil)
}

func (g *stocksGateway) GetItem(ctx context.Context, id string) (*pb.StockItem, error) {
	conn, err := discovery.ServiceConnection(context.Background(), stockServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewStockServiceClient(conn)

	return c.GetStockItem(ctx, &pb.GetStockItemRequest{Id: id})
}

func (g *stocksGateway) UpdateItem(ctx context.Context, id string, p *pb.UpdateStockItemRequest) (*pb.StockItem, error) {
	conn, err := discovery.ServiceConnection(context.Background(), stockServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewStockServiceClient(conn)

	p.Id = id

	return c.UpdateStockItem(ctx, p)
}

func (g *stocksGateway) UpdateStockQuantity(ctx context.Context, id string, quantity int) (*pb.StockItem, error) {
	conn, err := discovery.ServiceConnection(context.Background(), stockServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewStockServiceClient(conn)

	return c.UpdateStockQuantity(ctx, &pb.UpdateStockQuantityRequest{
		ID:       id,
		Quantity: int64(quantity),
	})
}

func (g *stocksGateway) DeleteItem(ctx context.Context, id string) error {
	conn, err := discovery.ServiceConnection(context.Background(), stockServiceName, g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	c := pb.NewStockServiceClient(conn)

	_, err = c.DeleteItem(ctx, &pb.DeleteItemRequest{ID: id})
	if err != nil {
		return err
	}

	return nil
}
