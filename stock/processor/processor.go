package processor

import (
	"time"

	pb "github.com/juxue97/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockProcessor interface {
	CreateProduct(*pb.Product) (string, string, error)
	UpdateProduct(prodID, priceID string, i Item) (string, error)
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
