package main

import (
	"context"
	"fmt"
	"time"

	"github.com/juxue97/common"
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
	mongoDB *mongo.Client
}

func NewStore(mongoDB *mongo.Client) *store {
	return &store{mongoDB: mongoDB}
}

func (s *store) GetItemsStock(ctx context.Context, ids []primitive.ObjectID) ([]*ItemStock, error) {
	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	// look for all the relevant document with the document ids
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find items: %v", err)
	}

	defer cursor.Close(ctx)

	var items []*ItemStock

	for cursor.Next(ctx) {
		var item ItemStock
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
	item.UpdatedAt = time.Now()
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
	update := bson.M{"$set": bson.M{
		"quantity":   quantity,
		"updated_at": time.Now(),
	}}
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

func (s *store) DeductStock(ctx context.Context, id string, quantity int) (*Item, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	var updatedItem Item
	update := bson.M{
		"$inc": bson.M{"quantity": -quantity}, // Deduct quantity
		"$set": bson.M{"updated_at": time.Now()},
	}
	filter := bson.M{
		"_id":      oID,
		"quantity": bson.M{"$gte": quantity}, // Ensure enough stock
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err = col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("insufficient stock or item not found")
		}
		return nil, fmt.Errorf("update failed: %v", err)
	}

	return &updatedItem, nil
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
	update := bson.M{"$set": bson.M{
		"active":     false,
		"updated_at": time.Now(),
	}}

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
		Active:      true,
		Currency:    item.Currency,
		Quantity:    item.Quantity,
		PriceID:     priceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if item.Metadata != nil {
		newProd.Metadata = item.Metadata
	}

	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	newProduct, err := col.InsertOne(ctx, newProd)
	id, ok := newProduct.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, common.ErrConvertID
	}

	return id, err
}
