package main

import (
	"context"

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

	itemsInStock, err := s.store.GetItems(ctx, itemIDs)
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

func (s *stockService) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	return s.store.GetItems(ctx, ids)
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
