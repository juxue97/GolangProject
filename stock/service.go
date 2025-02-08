package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"github.com/juxue97/stock/gateway"
	"github.com/juxue97/stock/processor"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type stockService struct {
	store           StockStore
	stripeProcessor processor.StockProcessor
	gateway         gateway.StocksGateway
}

func NewStockService(store StockStore, stripeProcessor processor.StockProcessor, gateway gateway.StocksGateway) *stockService {
	return &stockService{
		store:           store,
		stripeProcessor: stripeProcessor,
		gateway:         gateway,
	}
}

func (s *stockService) CheckIfItemInStock(ctx context.Context, p []*pb.ItemsWithQuantity) (bool, []*pb.Item, error) {
	itemIDs := make([]primitive.ObjectID, 0, len(p)) // Preallocate slice

	// Create a map for quick lookup of requested items
	requestedItems := make(map[string]int32)

	for _, item := range p {
		oID, err := primitive.ObjectIDFromHex(item.ID)
		if err != nil {
			return false, nil, fmt.Errorf("invalid item ID: %s", item.ID)
		}
		itemIDs = append(itemIDs, oID)
		requestedItems[item.ID] = item.Quantity
	}

	// Fetch items from stock
	itemsInStock, err := s.store.GetItemsStock(ctx, itemIDs)
	if err != nil {
		return false, nil, err
	}

	// If no items found, return false immediately
	if len(itemsInStock) == 0 {
		return false, nil, nil
	}

	// Map to store available items
	itemsInStockPB := make([]*pb.Item, 0, len(itemsInStock))

	// Check stock and prepare response
	for _, stockItem := range itemsInStock {
		stockID := stockItem.ID.Hex()
		requiredQty, exists := requestedItems[stockID]

		// If the requested item exists and quantity is insufficient, return false
		if exists && int32(stockItem.Quantity) < requiredQty {
			return false, nil, nil
		}

		// Add to response only if it exists in the request
		if exists {
			itemsInStockPB = append(itemsInStockPB, &pb.Item{
				ID:       stockID,
				Name:     stockItem.Name,
				Quantity: requiredQty,
				PriceID:  stockItem.PriceID,
				// Initialize other fields if needed
			})
		}
	}

	return true, itemsInStockPB, nil
}

func (s *stockService) GetItems(ctx context.Context) ([]*Item, error) {
	item, err := s.store.GetItems(ctx)
	if err != nil {
		return nil, err
	}
	if len(item) <= 0 {
		return nil, common.ErrNoDoc
	}

	return item, nil
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

func (s *stockService) DeductStock(ctx context.Context, id string, quantity int) (*Item, error) {
	return s.store.DeductStock(ctx, id, quantity)
}

func (s *stockService) UpdateStock(ctx context.Context, id string, quantity int) (*Item, error) {
	return s.store.UpdateStock(ctx, id, quantity)
}

func (s *stockService) GetOrderService(ctx context.Context, o *pb.Order) ([]*pb.Item, error) {
	res, err := s.gateway.GetOrder(ctx, o)
	if err != nil {
		return nil, err
	}

	return res.Items, nil
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
