package main

import (
	"context"
	"log"

	pb "github.com/juxue97/common/api"
)

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{store: store}
}

func (s *service) createOrder(ctx context.Context, payload *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	id, err := s.store.Create(ctx, payload, items)
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID:          id,
		CustomerID:  payload.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	}

	return o, nil
}

func (s *service) getOrder(ctx context.Context, payload *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, payload.OrderID, payload.CustomerID)
}

func (s *service) validateOrder(ctx context.Context, payload *pb.CreateOrderRequest) ([]*pb.Item, error) {
	mergedItems := mergeItemsQuantities(payload.Items)

	// validate with stock service
	log.Println(mergedItems)

	// temporary
	var itemsWithPrice []*pb.Item
	for _, item := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1QmgLeQoT0OvyN0AztngBGhM",
			Quantity: item.Quantity,
			ID:       item.ID,
		})
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
