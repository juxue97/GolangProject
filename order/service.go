package main

import (
	"context"
	"fmt"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/order/gateway"
)

type service struct {
	store   OrderStore
	gateway gateway.StocksGateway
}

func NewService(store OrderStore, gateway gateway.StocksGateway) *service {
	return &service{
		store:   store,
		gateway: gateway,
	}
}

func (s *service) createOrder(ctx context.Context, payload *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	id, err := s.store.Create(ctx, Order{
		CustomerID:  payload.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	})

	o := &pb.Order{
		ID:          id.Hex(),
		CustomerID:  payload.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	}

	return o, err
}

func (s *service) getOrder(ctx context.Context, payload *pb.GetOrderRequest) (*pb.Order, error) {
	o, err := s.store.Get(ctx, payload.OrderID, payload.CustomerID)
	if err != nil {
		return nil, err
	}
	return o.ToProto(), nil
}

func (s *service) validateOrder(ctx context.Context, payload *pb.CreateOrderRequest) ([]*pb.Item, error) {
	mergedItems := mergeItemsQuantities(payload.Items)

	// validate with stock service
	inStock, itemsWithPrice, err := s.gateway.CheckIfItemsInStock(ctx, payload.CustomerID, mergedItems)
	if err != nil {
		return nil, err
	}

	if !inStock {
		return nil, fmt.Errorf("insufficient stock amount")
	}

	return itemsWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	// algorithm to merge those items quantites based on customerID
	merged := make([]*pb.ItemsWithQuantity, 0)
	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				found = true
				finalItem.Quantity += item.Quantity
				break
			}
		}
		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}

func (s *service) updateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	if err := s.store.Update(ctx, o.ID, o); err != nil {
		return nil, err
	}

	return o, nil
}
