package main

import (
	"context"
	"log"

	pb "github.com/juxue97/common/api"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
}

func NewGRPCHandler(gRPCServer *grpc.Server, service OrderService) {
	handler := &gRPCHandler{service: service}
	pb.RegisterOrderServiceServer(gRPCServer, handler)
}

func (h *gRPCHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("Order Received!")

	h.service.validateOrder(ctx, payload)
	o := &pb.Order{
		ID: "42",
	}
	return o, nil
}
