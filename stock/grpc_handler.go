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

func (g *gRPCHandler) CreateItem(ctx context.Context, payload *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	oID, err := g.service.CreateItem(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &pb.CreateItemResponse{
		ObjectID: oID.Hex(),
	}, nil
}

func (g *gRPCHandler) GetStockItems(ctx context.Context, req *pb.Empty) (*pb.GetStockItemsResponse, error) {
	items, err := g.service.GetItems(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetStockItemsResponse{
		Items: items,
	}, nil
}

func (g *gRPCHandler) GetStockItem(ctx context.Context, p *pb.GetStockItemRequest) (*pb.StockItem, error) {
	item, err := g.service.GetItem(ctx, p.Id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (g *gRPCHandler) UpdateItem(ctx context.Context, p *pb.UpdateStockItemRequest) (*pb.StockItem, error) {
	item, err := g.service.UpdateItem(ctx, p.Id, p)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (g *gRPCHandler) UpdateStockQuantity(ctx context.Context, p *pb.UpdateStockQuantityRequest) (*pb.StockItem, error) {
	item, err := g.service.UpdateStock(ctx, p.ID, int(p.Quantity))
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (g *gRPCHandler) DeleteItem(ctx context.Context, p *pb.DeleteItemRequest) (*pb.Empty, error) {
	err := g.service.DeleteItem(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
