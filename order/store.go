package main

import (
	"context"
	"fmt"

	pb "github.com/juxue97/common/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName         = "orders"
	CollectionName = "orders"
)

type store struct {
	mongoDB *mongo.Client
}

func NewStore(mongoDB *mongo.Client) *store {
	return &store{mongoDB: mongoDB}
}

func (s *store) Create(ctx context.Context, o Order) (primitive.ObjectID, error) {
	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	newOrder, err := col.InsertOne(ctx, o)
	id, ok := newOrder.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to convert inserted ID to primitive.ObjectID")
	}

	return id, err
}

func (s *store) Get(ctx context.Context, orderID string, customerID string) (*Order, error) {
	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	oID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, err
	}

	var o Order
	filter := bson.M{
		"_id":        oID,
		"customerID": customerID,
	}

	err = col.FindOne(ctx, filter).Decode(&o)

	return &o, err
}

func (s *store) Update(ctx context.Context, id string, o *pb.Order) error {
	col := s.mongoDB.Database(DbName).Collection(CollectionName)

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": oID,
	}

	update := bson.M{
		"$set": bson.M{
			"paymentLink": o.PaymentLink,
			"status":      o.Status,
		},
	}

	result, err := col.UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with id %s", id)
	}
	return err
}
