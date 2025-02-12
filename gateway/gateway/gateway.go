package gateway

import (
	"context"

	pb "github.com/juxue97/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(context.Context, string, string) (*pb.Order, error)
}

type StocksGateway interface {
	CreateItem(ctx context.Context, p *pb.CreateItemRequest) (*pb.CreateItemResponse, error)
	GetItems(ctx context.Context) (*pb.GetStockItemsResponse, error)
	GetItem(ctx context.Context, id string) (*pb.StockItem, error)
	UpdateItem(ctx context.Context, id string, p *pb.UpdateStockItemRequest) (*pb.StockItem, error)
	UpdateStockQuantity(ctx context.Context, id string, quantity int) (*pb.StockItem, error)
	DeleteItem(ctx context.Context, id string) error
}
