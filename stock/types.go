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
	GetItems(context.Context) ([]*Item, error)
	GetItem(ctx context.Context, id string) (*Item, error)
	CreateItem(ctx context.Context, p *CreateItemRequest) (primitive.ObjectID, error)
	UpdateItem(ctx context.Context, id string) (*Item, error)
	UpdateStock(ctx context.Context, id string, quantity int) (*Item, error)
	DeleteItem(ctx context.Context, id string) error
	DeductStock(ctx context.Context, id string, quantity int) (*Item, error)
}

type StockStore interface {
	GetItemsStock(context.Context, []primitive.ObjectID) ([]*ItemStock, error)
	GetItems(ctx context.Context) ([]*Item, error)
	GetItem(context.Context, string) (*Item, error)
	CreateItem(ctx context.Context, prodID string, priceID string, item *CreateItemRequest) (primitive.ObjectID, error)
	UpdateItem(ctx context.Context, id string, item processor.Item) (*Item, error)
	UpdateStock(ctx context.Context, id string, quantity int) (*Item, error)
	DeleteItem(ctx context.Context, id string) error
	DeductStock(ctx context.Context, id string, quantity int) (*Item, error)
}

// can upload picture
type CreateItemRequest struct {
	Name        string            `json:"name" validate:"required,max=1000"`
	Description string            `json:"description" validate:"required,max=1000"`
	Price       float64           `json:"price" validate:"required"`
	Currency    string            `json:"currency" validate:"required"`
	Quantity    int64             `json:"quantity" validate:"required,min=1"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type CreateItemResponse struct {
	ObjectID primitive.ObjectID
}

type UpdateItemRequest struct {
	Name        string            `json:"name" validate:"omitempty,max=1000"`
	Description string            `json:"description" validate:"omitempty,max=1000"`
	Price       float64           `json:"price" validate:"omitempty"`
	Currency    string            `json:"currency" validate:"omitempty"`
	Quantity    int64             `json:"quantity" validate:"omitempty,min=1"`
	Metadata    map[string]string `json:"metadata" validate:"omitempty"`
	Active      bool              `json:"active" validate:"omitempty" default:"true"`
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
