package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"github.com/juxue97/stock/processor"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type stockService struct {
	store           StockStore
	stripeProcessor processor.StockProcessor
}

func NewStockService(store StockStore, stripeProcessor processor.StockProcessor) *stockService {
	return &stockService{
		store:           store,
		stripeProcessor: stripeProcessor,
	}
}

func (s *stockService) CheckIfItemInStock(ctx context.Context, p []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	itemIDs := make([]string, 0)

	for _, item := range p {
		itemIDs = append(itemIDs, item.ID)
	}

	itemsInStock, err := s.store.GetItemsStock(ctx, itemIDs)
	if err != nil {
		return false, nil, err
	}

	if len(itemsInStock) == 0 {
		return false, itemsInStock, err
	}

	for _, stockItem := range itemsInStock {
		for _, reqItem := range p {
			if stockItem.ID == reqItem.ID && stockItem.Quantity < reqItem.Quantity {
				return false, itemsInStock, nil
			}
		}
	}

	items := make([]*pb.Item, 0)
	for _, stockItem := range itemsInStock {
		for _, reqItem := range p {
			if stockItem.ID == reqItem.ID {
				items = append(items, &pb.Item{
					ID:       stockItem.ID,
					Name:     stockItem.Name,
					Quantity: stockItem.Quantity,
					PriceID:  stockItem.PriceID,
				})
			}
		}
	}

	return true, items, nil
}

func (s *stockService) GetItems(ctx context.Context) ([]*Item, error) {
	return s.store.GetItems(ctx)
}

func (s *stockService) GetItem(ctx context.Context, id string) (*Item, error) {
	return s.store.GetItem(ctx, id)
}

func (s *stockService) UpdateItem(ctx context.Context, id string, p *UpdateItemRequest) (*Item, error) {
	// Note: product id never change, but priceID will
	// first of all, obtain the product ID from MongoDB first
	oldItem, err := s.store.GetItem(ctx, id)
	if err != nil {
		return nil, err
	}
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload to JSON:%v", err)
	}

	var updateMap processor.Item
	if err := json.Unmarshal(jsonPayload, &updateMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON payload: %w", err)
	}

	var newPriceID string
	newPriceID, err = s.stripeProcessor.UpdateProduct(oldItem.ProductID, oldItem.PriceID, updateMap)
	if err != nil {
		return nil, err
	}
	if newPriceID != "" {
		updateMap.PriceID = newPriceID
	}

	return s.store.UpdateItem(ctx, id, updateMap)
}

func (s *stockService) UpdateStock(ctx context.Context, id string, quantity int) (*Item, error) {
	return s.store.UpdateStock(ctx, id, quantity)
}

func (s *stockService) DeleteItem(ctx context.Context, id string) error {
	oldItem, err := s.store.GetItem(ctx, id)
	if err != nil {
		return err
	}
	oldItem.Active = false
	param := &processor.Item{
		Active: false,
	}
	_, err = s.stripeProcessor.UpdateProduct(oldItem.ProductID, oldItem.PriceID, *param)
	if err != nil {
		return err
	}

	return s.store.DeleteItem(ctx, id)
}

func (s *stockService) CreateItem(ctx context.Context, p *CreateItemRequest) (primitive.ObjectID, error) {
	// call the stripe api to create new product
	// get the productID, Name and PriceID
	if p.Quantity <= 0 {
		return primitive.NilObjectID, common.ErrNoQuantity
	}

	product := &pb.Product{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Currency:    p.Currency,
		// TaxCode: item.TaxCode,
	}

	prodID, priceID, err := s.stripeProcessor.CreateProduct(product)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id, err := s.store.CreateItem(ctx, prodID, priceID, p)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return id, nil
}
