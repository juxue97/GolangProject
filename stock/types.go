package main

import (
	"context"
	"time"

	pb "github.com/juxue97/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockService interface {
	CheckIfItemInStock(context.Context, []*pb.ItemsWithQuantity) (bool, []*pb.Item, error)
	GetItems(context.Context, []string) ([]*pb.Item, error)
	CreateItem(ctx context.Context, p *CreateItemRequest) (primitive.ObjectID, error)
}

type StockStore interface {
	GetItems(context.Context, []string) ([]*pb.Item, error)
	GetItem(context.Context, string) (*pb.Item, error)
	CreateItem(ctx context.Context, prodID string, priceID string, item *CreateItemRequest) (primitive.ObjectID, error)
}

// can upload picture
type CreateItemRequest struct {
	Name        string  `json:"name" validate:"required,max=1000"`
	Description string  `json:"description" validate:"required,max=1000"`
	Price       float32 `json:"price" validate:"required"`
	Currency    string  `json:"currency" validate:"required"`
	Quantity    int64   `json:"quantity" validate:"required"`
}

type CreateItemResponse struct {
	ObjectID primitive.ObjectID
}

type Item struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProductID   string             `bson:"productID,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Price       float32            `bson:"price,omitempty"`
	Currency    string             `bson:"currency,omitempty"`
	Quantity    int64              `bson:"quantity,omitempty"`
	PriceID     string             `bson:"priceID,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}
