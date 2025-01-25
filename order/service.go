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

func (s *service) createOrder(ctx context.Context) error {
	return nil
}

func (s *service) validateOrder(ctx context.Context, payload *pb.CreateOrderRequest) error {
	mergedItems := mergeItemsQuantities(payload.Items)

	// validate with stock service
	log.Println(mergedItems)
	return nil
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
