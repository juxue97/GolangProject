package main

import (
	"context"
	"errors"

	pb "github.com/juxue97/common/api"
)

var orders = make([]*pb.Order, 0)

type store struct {
	// MongoDB
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, payload *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
	id := "42"

	orders = append(orders, &pb.Order{
		ID:         id,
		CustomerID: payload.CustomerID,
		Status:     "pending",
		Items:      items,
	})

	return id, nil
}

func (s *store) Get(ctx context.Context, orderID string, customerID string) (*pb.Order, error) {
	for _, o := range orders {
		if o.ID == orderID && o.CustomerID == customerID {
			return o, nil
		}
	}

	return nil, errors.New("order not found")
}

func (s *store) Update(ctx context.Context, id string, o *pb.Order) error {
	for i, order := range orders {
		if order.ID == id {
			orders[i].Status = o.Status
			orders[i].PaymentLink = o.PaymentLink
			return nil
		}
	}

	return nil
}
