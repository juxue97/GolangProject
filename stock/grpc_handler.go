package main

import (
	"context"

	pb "github.com/juxue97/common/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedStockServiceServer
	service *loggingMiddleware
	channel *amqp.Channel
}

func NewGRPCHandler(gRPCServer *grpc.Server, service *loggingMiddleware, channel *amqp.Channel) {
	handler := &gRPCHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterStockServiceServer(gRPCServer, handler)
}

func (g *gRPCHandler) CheckIfItemsInStock(ctx context.Context, payload *pb.CheckIfItemsInStockRequest) (*pb.CheckIfItemsInStockResponse, error) {
	inStock, items, err := g.service.CheckIfItemInStock(ctx, payload.Items)
	if err != nil {
		return nil, err
	}

	return &pb.CheckIfItemsInStockResponse{
		InStock: inStock,
		Items:   items,
	}, nil
}

func (g *gRPCHandler) GetItems(ctx context.Context, payload *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	items, err := g.service.GetItems(ctx, payload.ItemIDs)
	if err != nil {
		return nil, err
	}

	return &pb.GetItemsResponse{
		Items: items,
	}, nil
}
