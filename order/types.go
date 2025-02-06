package main

import (
	"context"

	pb "github.com/juxue97/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService interface {
	createOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	validateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	getOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	updateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type OrderStore interface {
	Create(context.Context, Order) (primitive.ObjectID, error)
	Get(context.Context, string, string) (*Order, error)
	Update(context.Context, string, *pb.Order) error
}

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID  string             `bson:"customerID,omitempty"`
	Status      string             `bson:"status,omitempty"`
	Items       []*pb.Item         `bson:"items,omitempty"`
	PaymentLink string             `bson:"paymentLink,omitempty"`
}

func (o *Order) ToProto() *pb.Order {
	return &pb.Order{
		ID:          o.ID.Hex(),
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
	}
}
