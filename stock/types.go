package main

import (
	"context"
	"time"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/stock/processor"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockService interface {
	CheckIfItemInStock(context.Context, []*pb.ItemsWithQuantity) (bool, []*pb.Item, error)
	GetOrderService(ctx context.Context, o *pb.Order) ([]*pb.Item, error)
	GetItems(context.Context) ([]*pb.StockItem, error)
	GetItem(ctx context.Context, id string) (*pb.StockItem, error)
	CreateItem(ctx context.Context, p *pb.CreateItemRequest) (primitive.ObjectID, error)
	UpdateItem(ctx context.Context, id string) (*pb.StockItem, error)
	UpdateStock(ctx context.Context, id string, quantity int) (*pb.StockItem, error)
	DeleteItem(ctx context.Context, id string) error
	DeductStock(ctx context.Context, id string, quantity int) (*Item, error)
}

type StockStore interface {
	GetItemsStock(context.Context, []primitive.ObjectID) ([]*ItemStock, error)
	GetItems(ctx context.Context) ([]*pb.StockItem, error)
	GetItem(context.Context, string) (*pb.StockItem, error)
	CreateItem(ctx context.Context, prodID string, priceID string, item *pb.CreateItemRequest) (primitive.ObjectID, error)
	UpdateItem(ctx context.Context, id string, item processor.Item) (*pb.StockItem, error)
	UpdateStock(ctx context.Context, id string, quantity int) (*pb.StockItem, error)
	DeleteItem(ctx context.Context, id string) error
	DeductStock(ctx context.Context, id string, quantity int) (*Item, error)
}

type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProductID   string             `bson:"productID,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Price       float64            `bson:"price,truncate,omitempty"`
	Currency    string             `bson:"currency,omitempty"`
	Quantity    int64              `bson:"quantity,omitempty"`
	Active      bool               `bson:"active"`
	PriceID     string             `bson:"priceID,omitempty"`
	Metadata    map[string]string  `bson:"metadata,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

type ItemStock struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Quantity int64              `bson:"quantity,omitempty"`
	PriceID  string             `bson:"priceID,omitempty"`
}
