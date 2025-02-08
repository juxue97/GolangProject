package main

import (
	"context"
	"fmt"
	"time"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"github.com/juxue97/stock/processor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *store) GetItemsStock(ctx context.Context, ids []string) ([]*pb.Item, error) {
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

func (s *store) GetItems(ctx context.Context) ([]*Item, error) {
	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []*Item
	for cursor.Next(ctx) {
		var item *Item
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("failed to decode item: %v", err)
		}
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return items, nil
}

func (s *store) GetItem(ctx context.Context, id string) (*Item, error) {
	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oID}

	var item *Item
	err = col.FindOne(ctx, filter).Decode(&item)
	if err == mongo.ErrNoDocuments {
		return nil, common.ErrNoDoc
	} else if err != nil {
		return nil, fmt.Errorf("failed to find item: %v", err)
	}

	return item, nil
}

func (s *store) UpdateItem(ctx context.Context, id string, item processor.Item) (*Item, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	var updatedItem *Item
	update := bson.M{"$set": item}
	filter := bson.M{"_id": oID}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.ErrNoDoc
		}
		return nil, fmt.Errorf("update failed: %v", err)
	}

	return updatedItem, nil
}

func (s *store) UpdateStock(ctx context.Context, id string, quantity int) (*Item, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	var updatedItem *Item
	update := bson.M{"$set": bson.M{"quantity": quantity}}
	filter := bson.M{"_id": oID}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.ErrNoDoc
		}
		return nil, fmt.Errorf("update failed: %v", err)
	}

	return updatedItem, nil
}

func (s *store) DeleteItem(ctx context.Context, id string) error {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	var item *Item
	filter := bson.M{
		"_id": oID,
	}
	update := bson.M{"$set": bson.M{"active": false}}

	err = col.FindOneAndUpdate(ctx, filter, update).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return common.ErrNoDoc
		}
		return fmt.Errorf("update failed: %v", err)
	}

	return nil
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
