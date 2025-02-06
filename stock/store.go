package main

import (
	"context"
	"fmt"
	"time"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName         = "stocks"
	CollectionName = "stocks"
)

type store struct {
	// mongoDB here
	stock   map[string]*pb.Item
	mongoDB *mongo.Client
}

func NewStore(mongoDB *mongo.Client) *store {
	return &store{stock: map[string]*pb.Item{
		"41": {
			ID:       "41",
			Name:     "ABC",
			PriceID:  "price_1QmgLeQoT0OvyN0AztngBGhM",
			Quantity: 10,
		},
		"42": {
			ID:       "42",
			Name:     "ABC",
			PriceID:  "price_1QmgLeQoT0OvyN0AztngBGhM",
			Quantity: 10,
		},
	}, mongoDB: mongoDB}
}

func (s *store) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	// look for all the relevant document with the document ids
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find items: %v", err)
	}

	defer cursor.Close(ctx)

	var items []*pb.Item
	for cursor.Next(ctx) {
		var item pb.Item
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("failed to decode item: %v", err)
		}
		items = append(items, &item)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return items, nil
}

func (s *store) GetItem(ctx context.Context, id string) (*pb.Item, error) {
	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	filter := bson.M{"_id": id}

	var item pb.Item
	err := col.FindOne(ctx, filter).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("item not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to find item: %v", err)
	}

	return &item, nil
}

func (s *store) CreateItem(ctx context.Context, prodID string, priceID string, item *CreateItemRequest) (primitive.ObjectID, error) {
	// create a new document
	newProd := Item{
		ProductID:   prodID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		Currency:    item.Currency,
		Quantity:    item.Quantity,
		PriceID:     priceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	newProduct, err := col.InsertOne(ctx, newProd)
	id, ok := newProduct.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, common.ErrConvertID
	}

	return id, err
}
